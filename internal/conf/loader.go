package conf

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/malivvan/vivid/internal/crypto"

	"github.com/BurntSushi/toml"
)

const (
	encryptLevel  = 0
	encryptPrefix = "$("
	encryptSuffix = ")"
	encryptTag    = "encrypt"
)

type File interface {
	io.Reader
	io.Writer
	io.Seeker
	Truncate(size int64) error
}

func LoadFile(path string, config interface{}) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	return Load(file, config)
}

func Load(file File, config interface{}) error {

	// 1. read bytes, if not exist write initial config
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	configBytes, err := io.ReadAll(file)
	if err != nil {
		var configBuf bytes.Buffer
		err := toml.NewEncoder(&configBuf).Encode(config)
		if err != nil {
			return err
		}
		configBytes = configBuf.Bytes()
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
		_, err = file.Write(configBytes)
		if err != nil {
			return err
		}
		err = file.Truncate(int64(len(configBytes)))
		if err != nil {
			return err
		}
	}

	// 2. decode config and encrypt all unencrypted strings
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = toml.Decode(string(configBytes), config)
	if err != nil {
		return err
	}
	if err := alterStrings(config, func(str string, secret string) (string, error) {
		if !(strings.HasPrefix(str, encryptPrefix) && strings.HasSuffix(str, encryptSuffix)) {
			encryptedData, err := crypto.Encrypt(crypto.ScryptLevel[encryptLevel], secret, []byte(str))
			if err != nil {
				return "", err
			}
			return encryptPrefix + base64.StdEncoding.EncodeToString(encryptedData) + encryptSuffix, nil
		}
		return str, nil
	}); err != nil {
		return err
	}

	// 3. if config was altered save to disk
	var newConfigBuf bytes.Buffer
	err = toml.NewEncoder(&newConfigBuf).Encode(config)
	if err != nil {
		return err
	}
	newConfigBytes := newConfigBuf.Bytes()
	if !bytes.Equal(newConfigBytes, configBytes) {
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
		_, err = file.Write(newConfigBytes)
		if err != nil {
			return err
		}
		err = file.Truncate(int64(len(newConfigBytes)))
		if err != nil {
			return err
		}
	}

	// 4. decrypt all strings
	if err := alterStrings(config, func(str string, secret string) (string, error) {
		if strings.HasPrefix(str, encryptPrefix) && strings.HasSuffix(str, encryptSuffix) {
			encryptedData, err := base64.StdEncoding.DecodeString(strings.TrimSuffix(strings.TrimPrefix(str, encryptPrefix), encryptSuffix))
			if err != nil {
				return "", err
			}
			decryptedData, _, err := crypto.Decrypt(secret, encryptedData)
			if err != nil {
				return "", err
			}
			return string(decryptedData), nil
		}
		return str, nil
	}); err != nil {
		return err
	}

	return nil
}

func alterStrings(v interface{}, f func(str string, secret string) (string, error)) error {
	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			for field.Kind() == reflect.Ptr {
				field = field.Elem()
			}
			switch field.Kind() {
			case reflect.Map:
				for _, key := range field.MapKeys() {
					err := alterStrings(field.MapIndex(key).Interface(), f)
					if err != nil {
						return err
					}
				}
			case reflect.Struct:
				err := alterStrings(val.Field(i).Addr().Interface(), f)
				if err != nil {
					return err
				}
			case reflect.String:
				if secret, ok := val.Type().Field(i).Tag.Lookup(encryptTag); ok {
					str, err := f(field.String(), secret)
					if err != nil {
						return err
					}
					field.SetString(str)
				}
			}
		}
	}
	return nil
}

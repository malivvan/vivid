package vivid

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kardianos/service"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {

	// Level of the logger. Valid options: debug, info, warn, error, disable.
	Level string `toml:"level"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `toml:"maxsize"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `toml:"maxage"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `toml:"maxbackups"`

	// Compress determines if the rotated log files should be compressed
	// using gzip.
	Compress bool `toml:"compress"`
}

type Logger struct {
	zerolog.Logger
	logfile *lumberjack.Logger
}

func (env *Environment) newLogger(name string) *Logger {
	logger := &Logger{}
	var writers = []io.Writer{}
	logConfig := reflect.ValueOf(env.config).Elem().FieldByName("Log").Addr().Interface().(*LogConfig)

	logDir := filepath.Join(env.appdir, "log")
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}

	logger.logfile = &lumberjack.Logger{
		Filename:   filepath.Join(logDir, name+".log"),
		MaxSize:    logConfig.MaxSize,
		MaxAge:     logConfig.MaxAge,
		MaxBackups: logConfig.MaxBackups,
		Compress:   logConfig.Compress,
	}
	writers = append(writers, logger.logfile)

	if service.Interactive() {
		writers = append(writers, zerolog.SyncWriter(zerolog.ConsoleWriter{Out: os.Stderr}))
	}
	logger.Logger = zerolog.New(io.MultiWriter(writers...)).With().Timestamp().Logger()
	return logger
}

func (logger *Logger) Close() (err error) {
	if logger.logfile != nil {
		err = logger.logfile.Close()
		logger.logfile = nil
	}
	return
}

type Watcher struct {
	waitgroup sync.WaitGroup
	running   atomic.Value
	modtime   time.Time
}

func newWatcher(path string, intervalMS int, callback func()) *Watcher {
	fsw := &Watcher{}

	// Get the initial modtime.
	fi, err := os.Stat(path)
	if err == nil {
		fsw.modtime = fi.ModTime()
	}

	// Start the file watcher routine.
	fsw.waitgroup.Add(1)
	fsw.running.Store(true)
	go func() {
		defer func() {
			fsw.running.Store(false)
			fsw.waitgroup.Done()
		}()
		for fsw.running.Load().(bool) {
			time.Sleep(time.Duration(intervalMS) * time.Millisecond)
			fi, err := os.Stat(path)
			if err != nil {
				continue
			}
			if fi.ModTime().After(fsw.modtime) {
				fsw.modtime = fi.ModTime()
				callback()
			}
		}
	}()

	return fsw
}

func (fsw *Watcher) Close() {
	fsw.running.Store(false)
	fsw.waitgroup.Wait()
}

package conf

import "github.com/malivvan/vlang/internal/plugin/influx"

type Config struct {
	Influx map[string]*influx.Config
}

package influx

import (
	"errors"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Config struct {
	Host   string `toml:"host" json:"host"`
	Token  string `toml:"token" json:"-" encrypt:"9rcWH9mkpE2ZUn3z8x2SnTcUHvpqBFYs"`
	Org    string `toml:"org" json:"org"`
	Bucket string `toml:"bucket" json:"bucket"`
}

var DefaultConfig = map[string]*Config{
	"default": {
		Host:   "http://localhost:8086",
		Token:  "my-token",
		Org:    "my-org",
		Bucket: "my-bucket",
	},
}

type PluginLoader struct {
	loader   map[string]func() *PluginInstance
	instance map[string]*PluginInstance
}

func (pl *PluginLoader) Name() string {
	return "Influx"
}

func (pl *PluginLoader) Version() string {
	return "0.0.1"
}

func (pl *PluginLoader) Func() func(string) interface{} {
	return func(name string) interface{} {
		return pl.loader[name]()
	}
}

func (pl *PluginLoader) Configure(v interface{}) error {
	configs, ok := v.(map[string]*Config)
	if !ok {
		return errors.New("invalid " + pl.Name() + " config")
	}
	for name := range configs {
		cfg := configs[name]
		pl.loader[name] = func() *PluginInstance {
			if plugin, ok := pl.instance[name]; ok {
				return plugin
			}
			instance := &PluginInstance{}
			instance.client = influxdb2.NewClient(cfg.Host, cfg.Token)
			instance.writeAPI = instance.client.WriteAPIBlocking(cfg.Org, cfg.Bucket)
			pl.instance[name] = instance
			return instance
		}
	}
	return nil
}

func (pl *PluginLoader) Cleanup() error {
	var closeErr error
	for _, plugin := range pl.instance {
		err := plugin.Destroy()
		if err != nil {
			closeErr = err
		}
	}
	pl.loader = make(map[string]func() *PluginInstance)
	pl.instance = make(map[string]*PluginInstance)
	return closeErr
}

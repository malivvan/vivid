package fyne

import (
	"errors"

	"fyne.io/fyne/v2"
)

type Config struct {
}

var DefaultConfig = map[string]*Config{
	"default": {},
}

type PluginLoader struct {
	App      fyne.App
	loader   map[string]func() *PluginInstance
	instance map[string]*PluginInstance
}

func (pl *PluginLoader) Name() string {
	return "Fyne"
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
		//cfg := configs[name]
		pl.loader[name] = func() *PluginInstance {
			if plugin, ok := pl.instance[name]; ok {
				return plugin
			}
			instance := &PluginInstance{app: pl.App}

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

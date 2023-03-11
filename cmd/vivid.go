package main

import (
	"github.com/malivvan/vivid"
	"github.com/malivvan/vivid/plugin/influx"
)

func main() {
	//	a := app.New()
	//w := a.NewWindow("Hello World")

	//w.SetContent(widget.NewLabel("Hello World!"))
	//w.Show()
	vivid.CLI(&struct {
		vivid.Config
		Influx map[string]*influx.Config `toml:"influx"`
	}{
		Config: vivid.DefaultConfig,
		Influx: influx.DefaultConfig,
	},
		&influx.PluginLoader{},
	)
	//a.Run()
}

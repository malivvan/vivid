package main

import (
	"github.com/malivvan/vivid"
	"github.com/malivvan/vivid/plugin/influx"
	"github.com/malivvan/vivid/updater"
)

func main() {
	//	a := app.New()
	//w := a.NewWindow("Hello World")

	//w.SetContent(widget.NewLabel("Hello World!"))
	//w.Show()

	rels, err := updater.GetReleases(vivid.AppRepo, vivid.AppBinary, false, false)
	if err != nil {
		panic(err)
	}
	for _, rel := range rels {
		println(rel.Tag())
	}

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

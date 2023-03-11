package fyne

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type PluginInstance struct {
	app fyne.App
}

func (pi *PluginInstance) Window() {
	w := pi.app.NewWindow("Hello World2")
	w.SetContent(widget.NewLabel("Hello World!2"))
	w.Show()

}

func (pi *PluginInstance) Destroy() error {
	return nil
}

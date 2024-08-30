package main

import (
	"os"

	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func main() {
	app := gtk.NewApplication("com.github.jwijenbergh.puregotk.gtk4.hello", gio.GApplicationFlagsNoneValue)
	// cleanup, no finalizers are used in this library
	defer app.Unref()
	// functions with callback arguments take function pointers
	// this is for internal re-use of callbacks
	actcb := func(_ gio.Application) {
		activate(app)
	}
	app.ConnectActivate(&actcb)

	if code := app.Run(len(os.Args), os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("Rokon: Control your Roku from your desktop")
	label := gtk.NewLabel("Hello from Go!")
	window.SetChild(&label.Widget)
	// cleanup, no finalizers are used in this library
	label.Unref()
	window.SetDefaultSize(500, 500)
	window.Present()
}

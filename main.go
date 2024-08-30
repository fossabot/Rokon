package main

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"

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
	err := sentry.Init(sentry.ClientOptions{
		Dsn:                "https://63c6c95f892988509925aaff62c839b3@o4504136997928960.ingest.us.sentry.io/4506838451945472",
		EnableTracing:      true,
		TracesSampleRate:   1.0,
		ProfilesSampleRate: 1.0,
		// Only enable Debug if the environment variable TRANSPARENT_TELEMETRY is set
		Debug: os.Getenv("TRANSPARENT_TELEMETRY") != "",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(2 * time.Second)

	if code := app.Run(len(os.Args), os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("Rokon: Control your Roku from your desktop")
	window.SetChild(&gtk.NewLabel("Hello from Go!").Widget)
	// cleanup, no finalizers are used in this library
	window.SetDefaultSize(400, 300)
	window.Present()
}

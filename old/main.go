package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"unsafe"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/getsentry/sentry-go"
)

/*
#cgo LDFLAGS: ./lib/liblibrokon_rust_sysinfo.a -ldl
#include "./lib/librokon_rust_sysinfo.h"
#include <stdlib.h>
*/
import "C"

func main() {
	fmt.Println("Starting Rokon. Now with more telemetry!")
	
	// C.sysinfo returns a cstring
	// we need to convert it to a go string
	// and then free the memory
	// https://stackoverflow.com/questions/40845677/how-to-convert-c-char-to-go-string
	sysinfo := C.sysinfo()
	defer C.free(unsafe.Pointer(sysinfo))
	err := sentry.Init(sentry.ClientOptions{
		Dsn:                "https://63c6c95f892988509925aaff62c839b3@o4504136997928960.ingest.us.sentry.io/4506838451945472",
		EnableTracing:      true,
		AttachStacktrace:   true,
		TracesSampleRate:   1.0,
		ProfilesSampleRate: 1.0,
		// Only enable Debug if the environment variable TRANSPARENT_TELEMETRY is set
		Debug: os.Getenv("TRANSPARENT_TELEMETRY") != "",
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// TRANSPARENT_TELEMETRY is set, so we can log the event and what data it's sending
			if os.Getenv("TRANSPARENT_TELEMETRY") != "" {
				fmt.Printf("Sending event: %s\n", event.Message)
				fmt.Printf("Event data: %v\n", event)
			}
			return event
		},
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	app := gtk.NewApplication("io.github.brycensranch.Rokon", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(2 * time.Second)

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	qualifier := func() string {
		switch {
		case os.Getenv("SNAP") != "":
			return " (Snap)"
		case os.Getenv("FLATPAK") != "":
			return " (Flatpak)"
		case os.Getenv("APPIMAGE") != "":
			return " (AppImage)"
		default:
			return ""
		}
	}()
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("Rokon: Control your Roku from your desktop")
	window.SetChild(&gtk.NewLabel("Hello from Go!").Widget)
	aboutWindow := gtk.NewAboutDialog()
	aboutWindow.SetProgramName("Rokon" + qualifier)
	aboutWindow.SetVersion(app.Version())
	aboutWindow.SetComments("Control your Roku TV from your desktop")
	aboutWindow.SetWebsite("https://github.com/BrycensRanch/Rokon")
	aboutWindow.SetWebsiteLabel("GitHub")
	aboutWindow.SetSystemInformation(
		C.GoString(C.sysinfo()) +
			"\n" + "GTK: " + strconv.Itoa(int(gtk.GetMajorVersion())) + "." + strconv.Itoa(int(gtk.GetMinorVersion())) + "." + strconv.Itoa(int(gtk.GetMicroVersion())))
	aboutWindow.SetCopyright("2024 Brycen G and contributors, but mostly Brycen")
	aboutWindow.SetWrapLicense(true)
	aboutWindow.SetModal(true)
	aboutWindow.SetDestroyWithParent(true)
	image := gtk.NewImageFromFile("assets/Rokon.png")
	// https://youtu.be/bLHL75H_VEM
	if image != nil {
		logo := image.Paintable()
		if logo == nil {
			fmt.Println("Warning: Could not load logo image from assets/Rokon.png")
			aboutWindow.SetLogoIconName("rokon")
		} else {
			aboutWindow.SetLogo(logo)
		}
	} else {
		fmt.Println("Warning: Could not load logo image from assets/Rokon.png")
		aboutWindow.SetLogoIconName("rokon")
	}

	// aboutWindow.SetAuthors([]string{"Brycen G. (BrycensRanch)"})
	aboutWindow.SetLicenseType(gtk.LicenseAGPL30)
	// window.SetChild(&aboutWindow.Window)

	aboutWindow.Present()
	aboutWindow.Focus()
	const windowSize = 400
	window.SetDefaultSize(windowSize, windowSize)
	// set window position to center

	window.Show()
}

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/getsentry/sentry-go"
)

func chooseNonEmpty(first, second string) string {
	if first != "" {
		return first
	}
	return second
}

func main() {
	fmt.Println("Starting Rokon. Now with more telemetry!")
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
				fmt.Printf("Sending event: %s\n", chooseNonEmpty(event.Type, event.Message))
				fmt.Printf("Event data: %v\n", event)
			}
			return event
		},
		// Integrations: [
		// 	sentry.Integration
		// ],
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	app := gtk.NewApplication("io.github.brycensranch.Rokon", gio.ApplicationDefaultFlags)
	if app.Version() == "" {
		app.SetVersion("0.0.0-SNAPSHOT")
	}
	app.ConnectActivate(func() { activate(app) })
	app.ConnectCommandLine(func(commandLine *gio.ApplicationCommandLine) int {
		return activateCommandLine(app, commandLine)
	})
	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(2 * time.Second)

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activateCommandLine(app *gtk.Application, commandLine *gio.ApplicationCommandLine) int {
	args := commandLine.Arguments() // Get the command-line arguments
	// Check if --version flag is present
	for _, arg := range args {
		if arg == "version" || arg == "--version" {
			// Print version info
			// commandLine.PrintLiteral("Now exiting")
			fmt.Println(applicationInfo(app))
			return 0 // Return 0 to indicate success
		}
	}
	commandLine.PrintLiteral("HI FROM COMMAND LINE RAHH")
	return 0 // or return another integer if needed
}

func applicationInfo(app *gtk.Application) string {
	qualifier := func() string {
		switch {
		case os.Getenv("SNAP") != "":
			return " (Snap)"
		case os.Getenv("FLATPAK") != "":
			return " (Flatpak)"
		case os.Getenv("APPIMAGE") != "":
			return " (AppImage)"
		case os.Getenv("CONTAINER") != "":
			return " (Container)"
		default:
			return ""
		}
	}()
	return fmt.Sprintf("Rokon %s%s", app.Version(), qualifier)
}

func activate(app *gtk.Application) {
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("Rokon: Control your Roku from your desktop")
	window.SetChild(&gtk.NewLabel("Hello from Go!").Widget)
	aboutWindow := gtk.NewAboutDialog()
	aboutWindow.SetProgramName(applicationInfo(app))
	aboutWindow.SetVersion(app.Version())
	aboutWindow.SetComments("Control your Roku TV from your desktop")
	aboutWindow.SetWebsite("https://github.com/BrycensRanch/Rokon")
	aboutWindow.SetWebsiteLabel("GitHub")
	aboutWindow.SetSystemInformation(
		("GTK: " + strconv.Itoa(int(gtk.GetMajorVersion())) + "." + strconv.Itoa(int(gtk.GetMinorVersion())) + "." + strconv.Itoa(int(gtk.GetMicroVersion()))))
	aboutWindow.SetCopyright("2024 Brycen G and contributors, but mostly Brycen")
	aboutWindow.SetWrapLicense(true)
	aboutWindow.SetModal(true)
	aboutWindow.SetDestroyWithParent(true)
	switch {
	case os.Getenv("SNAP") != "":
		image := gtk.NewImageFromFile(os.Getenv("SNAP") + "/meta/gui/icon.png")
		if image != nil {
			logo := image.Paintable()
			if logo != nil {
				aboutWindow.SetLogo(logo)
			}
		}
	case os.Getenv("FLATPAK") != "":
		image := gtk.NewImageFromFile("/app/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png")
		if image != nil {
			logo := image.Paintable()
			if logo != nil {
				aboutWindow.SetLogo(logo)
			}
		}
	default:
		// Assume native packaging
		aboutWindow.SetLogoIconName("io.github.brycensranch.Rokon")
		sentry.CaptureMessage("Something went wrong " + applicationInfo(app))

		if os.Getenv("CONTAINER") != "" {
			log.Println("Running in a container, the logo icon may not be displayed due to wrong path")
		}
	}
	// Capture an error and send it to Sentry
	// err := fmt.Errorf("something went wrong!")
	// sentry.CaptureException(err)
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

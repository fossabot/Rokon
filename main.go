package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/xdg"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/getsentry/sentry-go"
	"github.com/go-resty/resty/v2"
	"github.com/koron/go-ssdp"
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
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	switch runtime.GOOS {
	case "windows", "darwin":
		fmt.Println("Running on Windows or macOS.")
		// Use GLib to set the GTK_CSD environment variable for Client-Side Decorations
		glib.Setenv("GTK_CSD", "1", true)

		// Call your GTK-related functions here if needed

	default:
		fmt.Println("Running on Linux or other OS.")
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

// Search for Rokus asynchronously and return via channel
func searchForRokus() chan []ssdp.Service {
	resultChan := make(chan []ssdp.Service)

	go func() {
		defer close(resultChan)

		discoveredRokus, err := ssdp.Search("roku:ecp", 5, "")
		if err != nil {
			sentry.CaptureException(err)
			log.Println("Error discovering Rokus:", err)
			return
		}

		if discoveredRokus != nil {
			resultChan <- discoveredRokus // Send results back to the main thread
		} else {
			resultChan <- nil // No Rokus found, send nil
		}
	}()

	return resultChan
}

// Show the "About" window
func showAboutWindow(parent *gtk.ApplicationWindow, app *gtk.Application) {
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
		aboutWindow.SetLogoIconName("rokon")

		if os.Getenv("CONTAINER") != "" {
			log.Println("Running in a container, the logo icon may not be displayed due to wrong path")
		}
	}
	aboutWindow.SetAuthors([]string{"Brycen G. (BrycensRanch)"})
	aboutWindow.SetLicenseType(gtk.LicenseAGPL30)
	parent.SetChild(&aboutWindow.Window)

	aboutWindow.Present()
	aboutWindow.Focus()
}

// Create the main menu
func createMenu(window *gtk.ApplicationWindow, app *gtk.Application) *gio.Menu {
	menu := gio.NewMenu()

	// Create "Example" menu item
	exampleMenu := gio.NewMenuItem("Example", "example")
	exampleSubMenu := gio.NewMenu()

	// "About This App" menu item
	aboutMenuItem := gio.NewMenuItem("About This App", "about")
	aboutMenuItem.Connect("activate", func() {
		showAboutWindow(window, app)
	})
	exampleSubMenu.AppendItem(aboutMenuItem)

	// "Check For Updates" menu item
	updateMenuItem := gio.NewMenuItem("Check For Updates", "check-for-updates")
	updateMenuItem.Connect("activate", func() {
		fmt.Println("Checking for updates...")
		// You can integrate update logic here
	})
	exampleSubMenu.AppendItem(updateMenuItem)

	// "Quit" menu item
	quitMenuItem := gio.NewMenuItem("Quit", "quit")
	quitMenuItem.Connect("activate", func() {
		app.Quit()
	})
	exampleSubMenu.AppendItem(quitMenuItem)

	exampleMenu.SetSubmenu(exampleSubMenu)
	menu.AppendItem(exampleMenu)
	return menu
}

func fetchImageAsPaintable(url string) (string, error) {
	tempDir := filepath.Join(xdg.CacheHome, "rokon")
	client := resty.New()
	resp, err := client.SetOutputDirectory(tempDir).EnableTrace().R().
		// SetDebug(true).
		EnableTrace().
		SetOutput(filepath.Join(tempDir, "device-image.png")).
		Get(url)
	if err != nil {
		return "", err
	}
	// Check if the request was successful
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to get image: status code %d", resp.StatusCode())
	}
	imagePath := filepath.Join(tempDir, "device-image.png")
	println(imagePath)
	// image := gtk.NewImageFromFile(imagePath)

	return imagePath, nil
}

func activate(app *gtk.Application) {
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("Rokon: Control your Roku from your desktop")
	window.SetChild(&gtk.NewLabel("Searching for Rokus on your network...").Widget)
	const windowSize = 400
	window.SetDefaultSize(800, windowSize)
	// Start searching for Rokus when the app is activated
	rokuChan := searchForRokus()

	// Goroutine that waits for Roku discovery to finish
	go func() {
		discoveredRokus := <-rokuChan // Receive the result from the Roku discovery

		// Use glib.IdleAdd to ensure UI updates happen on the main thread
		glib.IdleAdd(func() {
			if discoveredRokus != nil {
				fmt.Println("Discovered Rokus:", discoveredRokus)
				window.SetChild(&gtk.NewLabel("Discovered Rokus:").Widget)
			} else {
				window.SetChild(&gtk.NewLabel("No Rokus discovered via SSDP!").Widget)
			}
		})

		// Once Roku discovery completes, run Resty logic
		if discoveredRokus != nil {
			client := resty.New()
			resp, err := client.R().
				EnableTrace().
				Get(discoveredRokus[0].Location + "/")

			if err != nil {
				fmt.Println("Error:", err)
				sentry.CaptureException(err)
			} else {
				fmt.Println("Trace Info:", resp.Request.TraceInfo())
				fmt.Println("Status Code:", resp.StatusCode())
				fmt.Println("Status:", resp.Status())
				fmt.Println("Proto:", resp.Proto())
				fmt.Println("Time:", resp.Time())
				fmt.Println("Received At:", resp.ReceivedAt())
				fmt.Println("Body:", resp)
			}
			notification := gio.NewNotification("Roku discovered")
			// Convert the list of Rokus into a single string
			var rokuList []string
			for _, roku := range discoveredRokus {
				// Assuming `USN` is the identifier you want to display for each Roku
				rokuList = append(rokuList, fmt.Sprintf("%s (%s)", roku.USN, roku.Location))
			}

			// Join the list into a single string, each on a new line
			rokuListString := strings.Join(rokuList, "\n")
			notification.SetBody(rokuListString)

			url := discoveredRokus[0].Location + "/device-image.png"
			// Create a new GIcon from the file
			imagePath, err := fetchImageAsPaintable(url)
			if err != nil {
				log.Println("Error getting image from URL:", err)
				return
			}
			// if bytesIcon.Icon == null {
			// 	log.Fatalln("bytesIcon is nil!")
			// }
			notification.SetIcon(gio.NewFileIcon(gio.NewFileForPath(imagePath)))
			// Set the icon for the notification
			notification.SetCategory("device")
			app.SendNotification("roku-discovered", notification)
		}
	}()

	window.SetVisible(true)
	// Create the main menu
	menu := createMenu(window, app)
	app.SetMenubar(menu)
	window.SetShowMenubar(true)

	keyController := gtk.NewEventControllerKey()
	keyController.SetName("keyController")
	window.AddController(keyController)

	keyController.Connect("key-pressed", func(controller *gtk.EventControllerKey, code uint) {
		println(code)
		if code == 93 { // Right-click
			println("right clicked")
		}
	})
	focusController := gtk.NewEventControllerFocus()
	focusController.SetName("focusController")
	window.AddController(focusController)
	gestureClick := gtk.NewGestureClick()
	gestureClick.SetName("gestureClick")
	gestureClick.Connect("pressed", func(_, numberOfPresses uint) {
		fmt.Println("Number of presses %s", numberOfPresses)
	})
	window.AddController(gestureClick)
	// window.Maximize()
}

// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"fmt"
	"log"
	"os"
<<<<<<< HEAD
<<<<<<< HEAD
	"path/filepath"
=======
>>>>>>> 82d6e8e (fix: add user analytics telemetry)
=======
	"path/filepath"
>>>>>>> 1a31384 (fix(telemetry): remove possible points of PII)
	"runtime"
	"strconv"
	"strings"
	"time"

<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/adrg/xdg"

=======
=======
	"github.com/adrg/xdg"

>>>>>>> cee6275 (feat: implemented roku discover and sync with master)
	"github.com/brycensranch/go-aptabase/pkg/aptabase/v1"
	"github.com/brycensranch/go-aptabase/pkg/osinfo/v1"
>>>>>>> 82d6e8e (fix: add user analytics telemetry)
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/getsentry/sentry-go"
	"github.com/go-resty/resty/v2"
	"github.com/koron/go-ssdp"
<<<<<<< HEAD
=======
	"encoding/xml"
>>>>>>> cee6275 (feat: implemented roku discover and sync with master)
)

var aptabaseClient *aptabase.Client // Package-level variable

// Root represents the root element of the XML.
type Root struct {
	XMLName    xml.Name  `xml:"root"`
	SpecVersion SpecVersion `xml:"specVersion"`
	Device      Device      `xml:"device"`
}

// SpecVersion holds the major and minor version numbers.
type SpecVersion struct {
	Major int `xml:"major"`
	Minor int `xml:"minor"`
}

// Device represents the device details.
type Device struct {
	DeviceType        string      `xml:"deviceType"`
	FriendlyName      string      `xml:"friendlyName"`
	Manufacturer      string      `xml:"manufacturer"`
	ManufacturerURL   string      `xml:"manufacturerURL"`
	ModelDescription  string      `xml:"modelDescription"`
	ModelName         string      `xml:"modelName"`
	ModelNumber       string      `xml:"modelNumber"`
	ModelURL          string      `xml:"modelURL"`
	SerialNumber      string      `xml:"serialNumber"`
	UDN               string      `xml:"UDN"`
	IconList          IconList    `xml:"iconList"`
	ServiceList       ServiceList  `xml:"serviceList"`
}

// IconList holds a list of icons.
type IconList struct {
	Icons []Icon `xml:"icon"`
}

// Icon represents an individual icon.
type Icon struct {
	MimeType string `xml:"mimetype"`
	Width    int    `xml:"width"`
	Height   int    `xml:"height"`
	Depth    int    `xml:"depth"`
	URL      string `xml:"url"`
}

// ServiceList holds a list of services.
type ServiceList struct {
	Services []Service `xml:"service"`
}

// Service represents an individual service.
type Service struct {
	ServiceType  string `xml:"serviceType"`
	ServiceID    string `xml:"serviceId"`
	ControlURL   string `xml:"controlURL"`
	EventSubURL  string `xml:"eventSubURL"`
	SCPDURL      string `xml:"SCPDURL"`
}

func chooseNonEmpty(first, second string) string {
	if first != "" {
		return first
	}
	return second
}

func getOSRelease() string {
	osName, osVersion := osinfo.GetOSInfo()
	return fmt.Sprintf("%s %s", osName, osVersion)
}

func createEvent(eventName string, eventData map[string]interface{}) aptabase.EventData {
	event := aptabase.EventData{
		EventName: eventName,
		Props:     eventData,
	}
	aptabaseClient.TrackEvent(event)
	return event
}

var (
	version              = "0.0.0-SNAPSHOT"
	isPackaged           = "false"
	packageFormat        = "native"
	telemetryOnByDefault = "true"
)

func main() {
	fmt.Println("Starting Rokon. Now with more telemetry!")
	err := sentry.Init(sentry.ClientOptions{
		Dsn:     "https://04484623ba4aa6cbb830e852178e9358@o4504136997928960.ingest.us.sentry.io/4507991443439616",
		Release: version,
		// Prevents HOSTNAME from being sent to Sentry.io (Is this PII? Anyway, don't need it)
		ServerName:         "",
		EnableTracing:      true,
		AttachStacktrace:   true,
		TracesSampleRate:   1.0,
		ProfilesSampleRate: 1.0,
		// Only enable Debug if the environment variable TRANSPARENT_TELEMETRY is set
		Debug: os.Getenv("TRANSPARENT_TELEMETRY") != "",
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// TRANSPARENT_TELEMETRY is set, so we can log the event and what data it's sending
			if os.Getenv("TRANSPARENT_TELEMETRY") != "" {
				log.Printf("Sending event: %s\n", chooseNonEmpty(event.Type, event.Message))
				log.Printf("Event ID: %v\n", chooseNonEmpty(hint.EventID, string(event.EventID)))
				log.Printf("Event data: %v\n", event)
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
	}
	aptabaseClient = aptabase.NewClient("A-US-0332858461", version, uint64(133), true, "")
	app := gtk.NewApplication("io.github.brycensranch.Rokon", gio.ApplicationDefaultFlags)
<<<<<<< HEAD
<<<<<<< HEAD
	switch runtime.GOOS {
	case "windows", "darwin":
		fmt.Println("Running on Windows or macOS.")
		// Use GLib to set the GTK_CSD environment variable for Client-Side Decorations
		glib.Setenv("GTK_CSD", "1", true)

		// Call your GTK-related functions here if needed

	default:
		fmt.Println("Running on Linux or other OS.")
	}
	if app.Version() == "" {
		app.SetVersion("0.0.0-SNAPSHOT")
=======
	if version == "" {
=======
	if version != "" {
>>>>>>> 1a31384 (fix(telemetry): remove possible points of PII)
		app.SetVersion(version)
	}
	switch runtime.GOOS {
	case "linux":
		release := getOSRelease()
		arch := runtime.GOARCH
		desktop := os.Getenv("XDG_CURRENT_DESKTOP")
		sessionType := os.Getenv("XDG_SESSION_TYPE")

		kdeSessionVersion := ""
		if desktop == "KDE" {
			kdeSessionVersion = os.Getenv("KDE_SESSION_VERSION")
		}

		log.Printf("Running on Linux. Specifically: %s %s with %s %s %s and %s\n",
			release, arch, desktop, os.Getenv("DESKTOP_SESSION"), kdeSessionVersion, sessionType)

		createEvent("linux_run", map[string]interface{}{
			"arch":           arch,
			"desktop":        desktop,
			"desktopVersion": kdeSessionVersion,
			"sessionType":    sessionType,
		})

		container := os.Getenv("container")
		if container != "" && container == "flatpak" {
			log.Println("Running from a Flatpak")
			createEvent("flatpak_run", map[string]interface{}{
				"flatpak":        container,
				"flatpakVersion": version, // Replace with your app version logic
			})
		} else if snap := os.Getenv("SNAP"); snap != "" {
			log.Println("Running from a Snap")
			createEvent("snap_run", map[string]interface{}{
				"snap":        snap,
				"snapVersion": version, // Replace with your app version logic
			})
		} else if appImage := os.Getenv("APPIMAGE"); appImage != "" {
			log.Println("Running from an AppImage")
			firejail := isRunningWithFirejail()

			if firejail {
				log.Println("Running from an AppImage with firejail")
				// Adjust telemetry or other settings as needed.
			}

			createEvent("appimage_run", map[string]interface{}{
				// The APPIAMGE variable often points to the user's /home/identifyingrealname/Applications
				// So strip that data out
				"appimage":           filepath.Base(appImage),
				"appimageVersion":    version, // Replace with your app version logic
				"firejail":           firejail,
				"desktopIntegration": os.Getenv("DESKTOPINTEGRATION"),
			})
		} else if isPackaged == "true" {
			log.Println("Running from a native package")
			createEvent("native_run", map[string]interface{}{
				"nativeVersion": version, // Replace with your app version logic
				// The PATH can sometimes reveal PII
				// "path":          os.Args[0],
				"packageFormat": packageFormat,
			})
		}
	case "windows":
		release := getOSRelease()
		arch := runtime.GOARCH
		log.Printf("Running on Windows. Specifically: %s %s\n",
			release, arch)

		if packageFormat == "portable" {
			log.Println("Running from a portable executable")
		}

		createEvent("windows_run", map[string]interface{}{
			"arch":          arch,
			"version":       version, // Replace with your app version logic
			"packageFormat": packageFormat,
		})
	case "darwin":
		release := getOSRelease()
		arch := runtime.GOARCH
		log.Printf("Running on macOS. Specifically: %s %s\n",
			release, arch)

		createEvent("macos_run", map[string]interface{}{
			"arch":    arch,
			"mas":     os.Getenv("MAS"),
			"version": version, // Replace with your app version logic
			"path":    os.Args[0],
		})
	default:
		log.Printf("Unsupported telemetry platform: %s %s %s. However, the application will continue.\n",
			runtime.GOOS, getOSRelease(), runtime.GOARCH)
		createEvent("unsupported_platform", map[string]interface{}{
			"platform": runtime.GOOS,
			"arch":     runtime.GOARCH,
			"version":  version, // Replace with your app version logic
			"path":     os.Args[0],
		})
>>>>>>> 82d6e8e (fix: add user analytics telemetry)
	}
	app.ConnectActivate(func() { activate(app) })
	app.ConnectCommandLine(func(commandLine *gio.ApplicationCommandLine) int {
		return activateCommandLine(app, commandLine)
	})
	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(2 * time.Second)
	aptabaseClient.Quit = true
	aptabaseClient.Stop()
	if code := app.Run(os.Args); code > 0 {
		sentry.Flush(2 * time.Second)
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
		case strings.Contains(version, "SNAPSHOT"):
			return " (Development)"
		default:
			return ""
		}
	}()
	return fmt.Sprintf("Rokon%s", qualifier)
}

<<<<<<< HEAD
=======
func isRunningWithFirejail() bool {
	appImage := os.Getenv("APPIMAGE")
	appDir := os.Getenv("APPDIR")
	return (appImage != "" && contains(appImage, "/run/firejail")) ||
		(appDir != "" && contains(appDir, "/run/firejail"))
}

// Helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

>>>>>>> cee6275 (feat: implemented roku discover and sync with master)
// Search for Rokus asynchronously and return via channel
func searchForRokus() chan []ssdp.Service {
	resultChan := make(chan []ssdp.Service)

	go func() {
		defer close(resultChan)

<<<<<<< HEAD
		discoveredRokus, err := ssdp.Search("roku:ecp", 5, "")
=======
		discoveredRokus, err := ssdp.Search("roku:ecp", 1, "")
>>>>>>> cee6275 (feat: implemented roku discover and sync with master)
		if err != nil {
			sentry.CaptureException(err)
			log.Println("Error discovering Rokus:", err)
			return
		}

		if discoveredRokus != nil {
<<<<<<< HEAD
			resultChan <- discoveredRokus // Send results back to the main thread
=======
			// Deduplicate based on LOCATION
			// Needed because the SSDP code runs on EVERY interface :)
			// So, if you have WiFi and Ethernet enabled, it will show two callbacks from your Roku TV.
			// The code is just that *good*
			locationMap := make(map[string]ssdp.Service)
			for _, roku := range discoveredRokus {
				locationMap[roku.Location] = roku
			}

			// Convert map back to a slice
			uniqueRokus := make([]ssdp.Service, 0, len(locationMap))
			for _, roku := range locationMap {
				uniqueRokus = append(uniqueRokus, roku)
			}
			resultChan <- uniqueRokus // Send results back to the main thread
>>>>>>> cee6275 (feat: implemented roku discover and sync with master)
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
	aboutWindow.SetModal(false)
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

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> cee6275 (feat: implemented roku discover and sync with master)
// Create the main menu
func createMenu(window *gtk.ApplicationWindow, app *gtk.Application) *gio.Menu {
	menu := gio.NewMenu()

	// Create "Example" menu item
<<<<<<< HEAD
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
=======
	exampleMenu := gio.NewMenuItem("Help", "example")
	exampleSubMenu := gio.NewMenu()

	// Add "About This App" action
	aboutAction := gio.NewSimpleAction("about", nil)
	aboutAction.Connect("activate", func() {
		showAboutWindow(window, app)
	})
	app.AddAction(aboutAction)
	aboutMenuItem := gio.NewMenuItem("About This App", "app.about")
	exampleSubMenu.AppendItem(aboutMenuItem)

	// Add "Check For Updates" action
	checkForUpdatesAction := gio.NewSimpleAction("check-for-updates", nil)
	checkForUpdatesAction.Connect("activate", func() {
		fmt.Println("Checking for updates...")
	})
	app.AddAction(checkForUpdatesAction)
	updateMenuItem := gio.NewMenuItem("Check For Updates", "app.check-for-updates")
	exampleSubMenu.AppendItem(updateMenuItem)

	// Add "Quit" action
	quitAction := gio.NewSimpleAction("quit", nil)
	quitAction.Connect("activate", func() {
		app.Quit()
	})
	app.AddAction(quitAction)
	quitMenuItem := gio.NewMenuItem("Quit", "app.quit")
>>>>>>> cee6275 (feat: implemented roku discover and sync with master)
	exampleSubMenu.AppendItem(quitMenuItem)

	exampleMenu.SetSubmenu(exampleSubMenu)
	menu.AppendItem(exampleMenu)
<<<<<<< HEAD
=======

>>>>>>> cee6275 (feat: implemented roku discover and sync with master)
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
<<<<<<< HEAD
	window.SetTitle("Rokon: Control your Roku from your desktop")
	window.SetChild(&gtk.NewLabel("Searching for Rokus on your network...").Widget)
	const windowSize = 400
=======
	// Create the main menu
	menu := createMenu(window, app)
	app.SetMenubar(menu)
	window.SetShowMenubar(true)
	window.SetTitle("Rokon: Control your Roku from your desktop")
	window.SetChild(&gtk.NewLabel("Searching for Rokus on your network...").Widget)
	windowSize := 400
>>>>>>> cee6275 (feat: implemented roku discover and sync with master)
	window.SetDefaultSize(800, windowSize)
	// Start searching for Rokus when the app is activated
	rokuChan := searchForRokus()

	// Goroutine that waits for Roku discovery to finish
	go func() {
		discoveredRokus := <-rokuChan // Receive the result from the Roku discovery

<<<<<<< HEAD
		// Use glib.IdleAdd to ensure UI updates happen on the main thread
		glib.IdleAdd(func() {
			if discoveredRokus != nil {
				fmt.Println("Discovered Rokus:", discoveredRokus)
				window.SetChild(&gtk.NewLabel("Discovered Rokus:").Widget)
			} else {
				window.SetChild(&gtk.NewLabel("No Rokus discovered via SSDP!").Widget)
			}
		})
=======
		// Perform the request and unmarshal directly into the Root struct
		var root Root
>>>>>>> cee6275 (feat: implemented roku discover and sync with master)

		// Once Roku discovery completes, run Resty logic
		if discoveredRokus != nil {
			client := resty.New()
			resp, err := client.R().
<<<<<<< HEAD
				EnableTrace().
=======
			SetResult(&root). // Set the result to automatically unmarshal the response
			EnableTrace().
>>>>>>> cee6275 (feat: implemented roku discover and sync with master)
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
<<<<<<< HEAD
			}
			notification := gio.NewNotification("Roku discovered")
			// Convert the list of Rokus into a single string
			var rokuList []string
			for _, roku := range discoveredRokus {
				// Assuming `USN` is the identifier you want to display for each Roku
				rokuList = append(rokuList, fmt.Sprintf("%s (%s)", roku.USN, roku.Location))
			}

			// Join the list into a single string, each on a new line
=======
						// Use glib.IdleAdd to ensure UI updates happen on the main thread
		glib.IdleAdd(func() {
			if discoveredRokus != nil {
				log.Println("Discovered Rokus:", discoveredRokus)
				// window.SetChild(&gtk.NewLabel("Discovered Rokus:").Widget)
					// Create a vertical box to hold the widgets
				vbox := gtk.NewBox(gtk.OrientationVertical, 5)
				window.SetChild(vbox)
							// Create a label with the friendly name and IP address
			labelText := fmt.Sprintf("Friendly Name: %s\nIP Address: %s",
				root.Device.FriendlyName, discoveredRokus[0].Location)
				label := gtk.NewLabel(labelText)
				vbox.Append(label) // Add label to the vertical box

			} else {
				window.SetChild(&gtk.NewLabel("No Rokus discovered via SSDP!").Widget)
			}
		})
			}
			notification := gio.NewNotification("Roku Discovered")

			// Limit verbose details like USN and Location, just list friendly names (e.g., device names or IPs).
			var rokuList []string
			for i, roku := range discoveredRokus {
				if i < 3 {
					// Include a simple label or name for each Roku. Replace `roku.USN` with a more user-friendly field if available.
					rokuList = append(rokuList, fmt.Sprintf("Roku Device %d: %v", i+1, roku.Location))
				}
			}

			// If more than 3 Rokus are discovered, add a note.
			if len(discoveredRokus) > 3 {
				rokuList = append(rokuList, fmt.Sprintf("...and %d more devices", len(discoveredRokus)-3))
			}

			// Join the list into a single string
>>>>>>> cee6275 (feat: implemented roku discover and sync with master)
			rokuListString := strings.Join(rokuList, "\n")
			notification.SetBody(rokuListString)

			url := discoveredRokus[0].Location + "/device-image.png"
			// Create a new GIcon from the file
			imagePath, err := fetchImageAsPaintable(url)
			if err != nil {
<<<<<<< HEAD
=======
				sentry.CaptureException(err)
>>>>>>> cee6275 (feat: implemented roku discover and sync with master)
				log.Println("Error getting image from URL:", err)
				return
			}
			// if bytesIcon.Icon == null {
			// 	log.Fatalln("bytesIcon is nil!")
			// }
			notification.SetIcon(gio.NewFileIcon(gio.NewFileForPath(imagePath)))
<<<<<<< HEAD
=======
			notification.SetDefaultAction("app.connect-roku")
>>>>>>> cee6275 (feat: implemented roku discover and sync with master)
			// Set the icon for the notification
			notification.SetCategory("device")
			app.SendNotification("roku-discovered", notification)
		}
	}()

	window.SetVisible(true)
<<<<<<< HEAD
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
=======
	window.SetVisible(true)
>>>>>>> e2d14b4 (refactor: use window set visible instead of deprecated method)
}
=======
	glib.IdleAdd(func() {
		log.Println("idleadd")
	})
>>>>>>> cee6275 (feat: implemented roku discover and sync with master)

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

// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/adrg/xdg"

	"github.com/brycensranch/go-aptabase/pkg/aptabase/v1"
	"github.com/brycensranch/go-aptabase/pkg/osinfo/v1"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/go-resty/resty/v2"
	"github.com/koron/go-ssdp"
)

var aptabaseClient *aptabase.Client // Package-level variable

// Root represents the root element of the XML.
type Root struct {
	XMLName     xml.Name    `xml:"root"`
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
	DeviceType       string      `xml:"deviceType"`
	FriendlyName     string      `xml:"friendlyName"`
	Manufacturer     string      `xml:"manufacturer"`
	ManufacturerURL  string      `xml:"manufacturerURL"`
	ModelDescription string      `xml:"modelDescription"`
	ModelName        string      `xml:"modelName"`
	ModelNumber      string      `xml:"modelNumber"`
	ModelURL         string      `xml:"modelURL"`
	SerialNumber     string      `xml:"serialNumber"`
	UDN              string      `xml:"UDN"`
	IconList         IconList    `xml:"iconList"`
	ServiceList      ServiceList `xml:"serviceList"`
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
	ServiceType string `xml:"serviceType"`
	ServiceID   string `xml:"serviceId"`
	ControlURL  string `xml:"controlURL"`
	EventSubURL string `xml:"eventSubURL"`
	SCPDURL     string `xml:"SCPDURL"`
}

// Structure to hold GitHub release information.
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

func getOSRelease() string {
	osName, osVersion := osinfo.GetOSInfo()
	return fmt.Sprintf("%s %s", osName, osVersion)
}

func createEvent(eventName string, eventData map[string]interface{}) {
	event := aptabase.EventData{
		EventName: eventName,
		Props:     eventData,
	}
	aptabaseClient.TrackEvent(event)
}

var (
	version              = "0.0.0-SNAPSHOT"
	isPackaged           = "false"
	packageFormat        = "native"
	telemetryOnByDefault = "true"
	commit               = "unknown"
)

func main() {
	fmt.Println("Starting Rokon. Now with more telemetry!")
	switch runtime.GOOS {
	case "windows", "darwin":
		fmt.Println("Running on Windows or macOS.")
		// Use GLib to set the GTK_CSD environment variable for Client-Side Decorations
		glib.Setenv("GTK_CSD", "1", true)
		os.Setenv("GTK_CSD", "1")
	default:
	}
	app := gtk.NewApplication("io.github.brycensranch.Rokon", gio.ApplicationDefaultFlags)
	aptabaseClient = aptabase.NewClient("A-US-0332858461", version, uint64(133), true, "")
	if version != "" {
		app.SetVersion(version)
	}
	fmt.Printf("Version %s commit %s", version, commit)
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
	}
	app.ConnectActivate(func() { activate(app) })
	app.ConnectCommandLine(func(commandLine *gio.ApplicationCommandLine) int {
		return activateCommandLine(app, commandLine)
	})
	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	aptabaseClient.Quit = true
	aptabaseClient.Stop()
	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activateCommandLine(_ *gtk.Application, commandLine *gio.ApplicationCommandLine) int {
	args := commandLine.Arguments() // Get the command-line arguments
	// Check if --version flag is present
	for _, arg := range args {
		if arg == "version" || arg == "--version" {
			// Print version info
			fmt.Println(applicationInfo())
			return 0 // Return 0 to indicate success
		}
	}
	commandLine.PrintLiteral("HI FROM COMMAND LINE RAHH")
	return 0
}

func applicationInfo() string {
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
	return "Rokon" + qualifier
}

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

// Search for Rokus asynchronously and return via channel.
func searchForRokus() chan []ssdp.Service {
	resultChan := make(chan []ssdp.Service)

	go func() {
		defer close(resultChan)

		discoveredRokus, err := ssdp.Search("roku:ecp", 1, "")
		if err != nil {
			log.Println("Error discovering Rokus:", err)
			return
		}

		if discoveredRokus != nil {
			resultChan <- discoveredRokus // Send results back to the main thread
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
		} else {
			resultChan <- nil // No Rokus found, send nil
		}
	}()

	return resultChan
}

// Show the "About" window.
func showAboutWindow(_ *gtk.ApplicationWindow, app *gtk.Application) {
	aboutWindow := gtk.NewAboutDialog()
	aboutWindow.SetProgramName(applicationInfo())
	aboutWindow.SetVersion(app.Version())
	aboutWindow.SetComments("Control your Roku TV from your desktop")
	aboutWindow.SetWebsite("https://github.com/BrycensRanch/Rokon")
	aboutWindow.SetWebsiteLabel("GitHub")
	//nolint:gosec // In GTK We trust.
	aboutWindow.SetSystemInformation(
		fmt.Sprintf("GTK: %d.%d.%d", int(gtk.GetMajorVersion()), int(gtk.GetMinorVersion()), int(gtk.GetMicroVersion())),
	)
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
	// aboutWindow.SetAuthors([]string{"Brycen G. (BrycensRanch)"})
	aboutWindow.SetLicenseType(gtk.LicenseAGPL30)

	aboutWindow.Present()
	aboutWindow.Focus()
}

// Create the main menu.
func createMenu(window *gtk.ApplicationWindow, app *gtk.Application) *gio.Menu {
	menu := gio.NewMenu()

	// Create "Example" menu item
	exampleMenu := gio.NewMenuItem("Example", "example")
	exampleSubMenu := gio.NewMenu()

	// "About This App" menu item
	aboutMenuItem := gio.NewMenuItem("About This App", "app.about")
	aboutMenuItem.Connect("activate", func() {
		showAboutWindow(window, app)
	})

	// "Check For Updates" menu item
	updateMenuItem := gio.NewMenuItem("Check For Updates", "app.check-for-updates")

	// "Quit" menu item
	quitMenuItem := gio.NewMenuItem("Quit", "quit")
	quitMenuItem.Connect("activate", func() {
		app.Quit()
	})

	aboutAction := gio.NewSimpleAction("about", nil)
	aboutAction.Connect("activate", func() {
		showAboutWindow(window, app)
	})
	app.AddAction(aboutAction)
	exampleSubMenu.AppendItem(aboutMenuItem)

	// Add "Check For Updates" action
	checkForUpdatesAction := gio.NewSimpleAction("check-for-updates", nil)
	checkForUpdatesAction.Connect("activate", func() {
		const url = "https://api.github.com/repos/BrycensRanch/Rokon/releases/latest"

		// Create a Resty client
		client := resty.New()

		// Make the request
		resp, err := client.R().SetResult(&GitHubRelease{}).Get(url)
		if err != nil {
			showDialog("Update Check Failed", fmt.Sprintf("Error fetching release information: %v", err), app)
			return
		}

		if resp.StatusCode() != 200 {
			showDialog("Update Check Failed", "Unable to fetch release info: "+resp.Status(), app)
			return
		}

		// Decode the response into the GitHubRelease struct
		release := resp.Result().(*GitHubRelease)

		// Compare versions
		if release.TagName != app.Version() {
			showDialog("Update Available", "A new version is available: "+release.TagName, app)
		} else {
			showDialog("Up to Date", "Your application is up to date.", app)
		}
	})

	app.AddAction(checkForUpdatesAction)
	exampleSubMenu.AppendItem(updateMenuItem)

	// Add "Quit" action
	quitAction := gio.NewSimpleAction("quit", nil)
	quitAction.Connect("activate", func() {
		app.Quit()
	})
	app.AddAction(quitAction)
	exampleSubMenu.AppendItem(quitMenuItem)

	exampleMenu.SetSubmenu(exampleSubMenu)
	menu.AppendItem(exampleMenu)

	return menu
}

// Function to show a dialog with the specified title and message.
func showDialog(title, message string, app *gtk.Application) {
	theWindow := gtk.NewWindow()
	dialog := gtk.NewMessageDialog(
		theWindow,
		gtk.DialogDestroyWithParent,
		gtk.MessageInfo,
		gtk.ButtonsNone,
	)

	dialog.SetTitle(title)
	dialog.SetApplication(app)
	dialog.SetModal(true)
	dialog.SetChild(gtk.NewLabel(message))
	theWindow.Show()
	dialog.Show()
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
	successfulHTTPCode := 200
	if resp.StatusCode() != successfulHTTPCode {
		return "", fmt.Errorf("failed to get image: status code %d", resp.StatusCode())
	}
	imagePath := filepath.Join(tempDir, "device-image.png")
	println(imagePath)
	// image := gtk.NewImageFromFile(imagePath)

	return imagePath, nil
}

func activate(app *gtk.Application) {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error fetching network interfaces:", err)
		return
	}

	for _, iface := range interfaces {
		// Get the interface status
		status := "down"
		if iface.Flags&net.FlagUp != 0 {
			status = "up"
		}

		// Determine the type of the interface
		var ifaceType string
		switch {
		case iface.Flags&net.FlagLoopback != 0:
			ifaceType = "loopback"
		case strings.Contains(iface.Name, "en") || strings.Contains(iface.Name, "eth"):
			ifaceType = "Ethernet"
		case strings.Contains(iface.Name, "wl"):
			ifaceType = "Wi-Fi"
		default:
			ifaceType = "Unknown"
		}

		// Print interface details
		fmt.Printf("Interface: %s, Status: %s, Type: %s\n", iface.Name, status, ifaceType)
	}

	window := gtk.NewApplicationWindow(app)
	window.SetTitle("Rokon: Control your Roku from your desktop")
	window.SetChild(&gtk.NewLabel("Searching for Rokus on your network...").Widget)
	// Create the main menu
	menu := createMenu(window, app)
	app.SetMenubar(menu)
	window.SetShowMenubar(true)
	window.SetTitle("Rokon: Control your Roku from your desktop")
	window.SetChild(&gtk.NewLabel("Searching for Rokus on your network...").Widget)
	windowWidth := 800
	windowHeight := 400
	window.SetDefaultSize(windowWidth, windowHeight)
	window.SetVisible(true)

	// Event controller setup
	keyController := gtk.NewEventControllerKey()
	keyController.SetName("keyController")
	window.AddController(keyController)

	keyController.Connect("key-pressed", func(controller *gtk.EventControllerKey, code uint) {
		println(controller.Name() + " " + strconv.FormatUint(uint64(code), 10))
		const (
			RightClickCode = uint(93) // Code representing a right-click
		)
		if code == RightClickCode {
			println("right clicked")
		}
	})

	focusController := gtk.NewEventControllerFocus()
	focusController.SetName("focusController")

	focusController.Connect("enter", func() {
		println("Keyboard focus entered!")
	})
	window.AddController(focusController)

	gestureClick := gtk.NewGestureClick()
	gestureClick.SetName("gestureClick")
	gestureClick.Connect("pressed", func(_, numberOfPresses uint) {
		log.Printf("Number of presses %v", numberOfPresses)
	})
	window.AddController(gestureClick)

	// window.Maximize()
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

		// Perform the request and unmarshal directly into the Root struct
		var root Root

		// Once Roku discovery completes, run Resty logic
		if discoveredRokus != nil {
			client := resty.New()
			resp, err := client.R().
				SetResult(&root). // Set the result to automatically unmarshal the response
				Get(discoveredRokus[0].Location + "/")

			if err != nil {
				fmt.Println("Error:", err)
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
			var rokuList []string
			const (
				MaxValue = int(3) // Maximum allowed Rokus to display in notification
			)
			for i, roku := range discoveredRokus {
				if i < MaxValue {
					rokuList = append(rokuList, fmt.Sprintf("Roku Device %d: %v", i+1, roku.Location))
				}
			}
			if len(discoveredRokus) > MaxValue {
				rokuList = append(rokuList, fmt.Sprintf("...and %d more devices", len(discoveredRokus)-MaxValue))
			}
			rokuListString := strings.Join(rokuList, "\n")
			notification.SetBody(rokuListString)

			url := discoveredRokus[0].Location + "/device-image.png"
			imagePath, err := fetchImageAsPaintable(url)
			if err != nil {
				log.Println("Error getting image from URL:", err)
				return
			}
			notification.SetIcon(gio.NewFileIcon(gio.NewFileForPath(imagePath)))
			notification.SetDefaultAction("app.connect-roku")
			notification.SetCategory("device")
			app.SendNotification("roku-discovered", notification)

			// UI update for discovered Rokus
			glib.IdleAdd(func() {
				if discoveredRokus != nil {
					log.Println("Discovered Rokus:", discoveredRokus)
					const spacing = int(5)
					vbox := gtk.NewBox(gtk.OrientationVertical, spacing)
					window.SetChild(vbox)
					labelText := fmt.Sprintf("Friendly Name: %s\nIP Address: %s",
						root.Device.FriendlyName, discoveredRokus[0].Location)
					label := gtk.NewLabel(labelText)
					vbox.Append(label) // Add label to the vertical box
				} else {
					window.SetChild(&gtk.NewLabel("No Rokus discovered via SSDP!").Widget)
				}
			})
		}
	}()
}

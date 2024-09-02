package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// This was tested on Wine staging 9.16 and Windows 10 and 11. It works on both.
// Linux cannot build this normally, as it's using Windows-specific functions
// To build this you must do `GOOS=windows go build`

func main() {
	// Load the ntdll.dll module
	hntdll, err := windows.LoadLibrary("ntdll.dll")
	if err != nil {
		fmt.Println("Not running on NT.")
		return
	}
	defer windows.FreeLibrary(hntdll)

	// Get the addresses of the Wine functions
	wineGetVersionProc, err := windows.GetProcAddress(hntdll, "wine_get_version")
	if err != nil {
		fmt.Println("Running under Windows.")
		return
	}

	wineGetHostVersionProc, err := windows.GetProcAddress(hntdll, "wine_get_host_version")
	if err != nil {
		fmt.Println("Failed to get wine_get_host_version function.")
		return
	}

	// Call wine_get_version
	versionPtr, _, _ := syscall.Syscall(wineGetVersionProc, 0, 0, 0, 0)
	version := windows.BytePtrToString((*byte)(unsafe.Pointer(versionPtr)))

	// Call wine_get_host_version
	var sysname, release uintptr
	/// Ignore what VSCode says about the number of arguments, this is correct. 
	syscall.Syscall(wineGetHostVersionProc, 2, uintptr(unsafe.Pointer(&sysname)), uintptr(unsafe.Pointer(&release)), 0)

	sysnameStr := windows.BytePtrToString((*byte)(unsafe.Pointer(sysname)))
	releaseStr := windows.BytePtrToString((*byte)(unsafe.Pointer(release)))

	fmt.Printf("Running under Wine %s under %s %s.\n", version, sysnameStr, releaseStr)
	// Now error out, to let the caller know
	os.Exit(1)
}

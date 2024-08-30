package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

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
	syscall.Syscall(wineGetHostVersionProc, 2, uintptr(unsafe.Pointer(&sysname)), uintptr(unsafe.Pointer(&release)), 0)

	sysnameStr := windows.BytePtrToString((*byte)(unsafe.Pointer(sysname)))
	releaseStr := windows.BytePtrToString((*byte)(unsafe.Pointer(release)))

	fmt.Printf("Running Wine %s under %s %s.\n", version, sysnameStr, releaseStr)
}

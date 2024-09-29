# Building Guide

Looking to contribute? Don't forget to check [CONTRIBUTING.md](./CONTRIBUTING.md).

## Prerequisites

Before you start contributing, please make sure you have the following:

- [Go](https://golang.org) installed for building, not needed at runtime!
- [GTK4](https://www.gtk.org) installed (don't forget gobject-introspection-1.0)
- [Git](https://git-scm.com) installed and configured with your GitHub account and your commits signed
- [Roku device](https://www.roku.com/products/roku-tv) (for local testing, not required for building)
- Building the application for Windows requires [MSYS2](https://www.msys2.org/) with the dependencies for GTK4, Golang  installed.
- Building the application for Windows from Linux/MacOS requires [cross-compiling](https://github.com/diamondburned/gotk4/issues/147) which is not supported by this project. There's libraries missing from Fedora's repositories that are required for cross-compiling GTK4 Golang applications. It's best to let the CI/CD pipeline handle the Windows builds. Commit and push your changes to the repository and let the CI/CD pipeline handle the rest.
- Patience and a willingness to learn

## Development Setup

> Keep in mind that building the project took 12 minutes on my blazingly fast i7-11800H processor with 16GB of RAM. It will take longer on slower hardware. Good luck!

### 🪟 Windows

#### Traditional processors (AMD64 or x64 architecture)

If you want to build the application on Windows, you will need to install [MSYS2](https://www.msys2.org/) and the dependencies for GTK4, Golang.

This can be accomplished with:

```powershell
winget install --id=MSYS2.MSYS2 -e
```

Then, open the MSYS2 CLANG64 terminal and run the following commands:

```bash
pacman -Syu
pacman -S --noconfirm git mingw-w64-clang-x86_64-go mingw-w64-clang-x86_64-gtk4 mingw-w64-clang-x86_64-gobject-introspection mingw-w64-clang-x86_64-gdb mingw-w64-clang-x86_64-toolchain make
```

#### Windows on ARM (like Snapdragon X Elite) (aarch64 or ARM64)

Windows 11 is **required** for Windows on ARM. The results binaries from `make` will be native.

```powershell
winget install --id=MSYS2.MSYS2 -e
```

Then, open the MSYS2 CLANGARM64 terminal and run the following commands:

```bash
pacman -Syu
pacman -S --noconfirm git mingw-w64-clang-aarch64-go mingw-w64-clang-aarch64-gtk4 mingw-w64-clang-aarch64-gobject-introspection mingw-w64-clang-aarch64-toolchain make
```

### 🍎 macOS

> macOS: `brew install git go gtk4 gobject-introspection make`

Note: You will need to install [Homebrew](https://brew.sh) to install the dependencies on macOS.

After that, follow the Linux build instructions, as they are the same for macOS (UNIX-like).

### 🐧 Linux

> You cannot build the application on Debian Stable, Ubuntu 22.04, PopOS, Linux Mint  or older distributions due to packages being too old in the software repositories. This applies to RHEL, OpenSUSE Leap, Mageia 9, and CentOS as well. You will need to use a newer distribution like Fedora, Ubuntu 24.04, Linux Mint 22, or OpenSUSE Tumbleweed/Slowroll.

> If you don't want to use a newer distribution, you can use a containerized build and runtime environment. This gets around your ancient software repositories. For this, you'll need [Distrobox](https://wiki.archlinux.org/title/Distrobox) (container) and choose the Fedora 40 image and then follow the instructions for Fedora.

> If you don't want to use Distrobox, you can always use our prebuilt AppImage which works on old and new distributions. Or our Snap package or Flatpak package.

> Fedora Linux 40+: `sudo dnf in -y git golang gtk4-devel @development-tools`

> Ubuntu 24.04+:

 ```bash
sudo apt install -y build-essential git libgtk-4-dev libgirepository1.0-dev software-properties-common make
# If using Rhino Linux, do not add the PPA. You already have a modern version of Golang 100% of the time.
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install -y golang-go
```

> Rhino Linux:

```bash
sudo nala install -y build-essential git libgtk-4-dev golang-go libgirepository1.0-dev make
```


> Alpine Linux:

```bash
sudo setup-apkrepos -c -1
sudo apk add --no-cache alpine-sdk gtk4.0-dev gobject-introspection-dev go make
```

> Arch Linux: `sudo pacman -S --noconfirm git base-devel go gtk4 gobject-introspection`

> OpenSUSE Tumbleweed/Slowroll: `sudo zypper in -y git go gtk4-devel gobject-introspection-devel make`

## Cloning and building the project

```bash
git clone https://github.com/BrycensRanch/Rokon
cd Rokon
# If your internet is slow, this WILL take awhile.
go mod download all
# This may take a while, CGO is slow.
make build
# If on Windows, do not add "sudo"
# On Windows & macOS make install does not natively integrate with your operating system. 
# All it does is add rokon to your $PATH.
# On most Linux distributions, it actually *properly* installs as any normal application as you'd expect.
sudo make install
```

Run the app in development mode:

```bash
wgo go run -v .
```

Note: To truly test the development version of the app, you will need a Roku device.

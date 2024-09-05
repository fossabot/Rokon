# Contributing to Rokon

Thank you for your interest in contributing to Rokon! We appreciate your help in making our Roku remote app even better.

## Prerequisites

Before you start contributing, please make sure you have the following:

- [Go](https://golang.org) installed
- [GTK4](https://www.gtk.org) installed
- [Git](https://git-scm.com) installed and configured with your GitHub account and your commits signed
- [Roku device](https://www.roku.com/products/roku-tv) (for local testing, not required for building)
- Builiding the application for Windows requires [MSYS2](https://www.msys2.org/) with the dependencies for GTK4, Golang, and Rust (with cargo) installed.
- Building the application for Windows from Linux/MacOS requires [cross-compiling](https://github.com/diamondburned/gotk4/issues/147) which is not supported by the project. That's not to say it's impossible, but it's not recommended. There's libraries missing from Fedora's repositories that are required for cross-compiling GTK4 Golang applications. It's best to let the CI/CD pipeline handle the Windows builds. Commit and push your changes to the repository and let the CI/CD pipeline handle the rest.
- Patience and a willingness to learn

## Development Setup

> Installing the dependencies depends on your operating system. Here, we have the most common operating systems and their package managers.

> Keep in mind that building the project took 12 minutes on my blazingly fast i7-11800H processor with 16GB of RAM. It will take longer on slower hardware. Good luck!

### 🪟 Windows

#### AMD64 ONLY

This application was developed on Linux. While it was built on Linux, you can use it on Windows. However, building the application on Windows is not recommended. If you want to build the application on Windows, you will need to install MSYS2 and the dependencies for GTK4, Golang

This can be accomplished with:

```powershell
winget install --id=MSYS2.MSYS2 -e

```

Then, open the MSYS2 MINGW64 terminal and run the following commands:

```bash
pacman -Syu
pacman -S --noconfirm git mingw-w64-x86_64-go mingw-w64-x86_64-gtk4 mingw-w64-x86_64-gobject-introspection
```

#### Windows on ARM (laughably untested, glhf)

If you're using Windows on ARM, you will need to do something a bit different. Make sure you're using Windows 11 and have Docker setup and installed. You can use the following powershell script (.ps1) to build the application for ARM64 on Windows:

```powershell
# filename: build.ps1
# Define Go version and build flags
$GO_VERSION = "latest"  # or "1.22" for bare minimum Go version
$GO_BUILD_FLAGS = "-x"

# Get GOCACHE and GOPATH using Go commands
$GOCACHE = go env GOCACHE
$GOPATH = go env GOPATH

function docker_build {
    param (
        [string]$GOARCH
    )

    Write-Host ":: Building for $GOARCH in MinGW container..."

    # Docker command to run Go build inside the MinGW container
    docker run --rm -it `
        -e GOCACHE=/go/.cache `
        -e GOARCH=$GOARCH `
        -v "$GOCACHE:/go/.cache" `
        -v "$GOPATH/src:/go/src" `
        -v "$GOPATH/pkg:/go/pkg" `
        -v "$(Get-Location):/go/work" `
        -w /go/work `
        "x1unix/go-mingw:$GO_VERSION" `
        go build $GO_BUILD_FLAGS -v -o "rokon-$GOARCH.exe" .
}

# Example usage of the docker_build function
docker_build "arm64"  # or any architecture like "386", "amd64", etc.
```

After you've done that, you're done. Do not read any further, go have a coffee or whatever Windows users do.

### 🍎 macOS (ultra proprietary garbage)

> macOS: `brew install git go gtk4 gobject-introspection`

Note: You will need to install [Homebrew](https://brew.sh) to install the dependencies on macOS.

After that, follow the Linux build instructions, as they are the same for macOS (UNIX-like).

### 🐧 Linux

> You cannot build the application on Debian Stable, Ubuntu 22.04, PopOS, Linux Mint  or older distributions due to packages being too old in the software repositories. This applies to RHEL, OpenSUSE Leap, Mageia,and CentOS as well. You will need to use a newer distribution like Fedora, Ubuntu 24.04, Linux Mint 22, or OpenSUSE Tumbleweed/Slowroll.

> If you don't wnat to use a newer distribution, you can use a containerized build and runtime environment. This gets around your ancient software repositories. For this, you'll need [Distrobox](https://wiki.archlinux.org/title/Distrobox) (container) and choose the Fedora 40 image and then follow the instructions for Fedora.

> If you don't want to use Distrobox, you can always use our prebuilt AppImage which works on old and new distributions. Or our Snap package or Flatpak package.

> Fedora Linux 40+: `sudo dnf in -y git golang gtk4-devel @development-tools`

> Ubuntu 24.04+:

 ```bash
sudo apt install -y build-essential git libgtk-4-dev libgirepository1.0-dev software-properties-common
# If using Rhino Linux, do not add the PPA. Rather just install the go-bin package with pacstall. This PPA is added for newer versions of Go for Ubuntu Stable releases.
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install -y golang-go
```

> Alpine Linux:

```bash
sudo setup-apkrepos -c -1
sudo apk add --no-cache alpine-sdk gtk4.0-dev gobject-introspection-dev go
```

> Arch Linux: `sudo pacman -S --noconfirm git base-devel go gtk4 gobject-introspection`

> OpenSUSE Tumbleweed/Slowroll: `sudo zypper in -y git go gtk4-devel gobject-introspection-devel`

## Cloning and building the project

```bash
git clone https://github.com/BrycensRanch/Rokon
cd Rokon
# If your internet is slow, this WILL take awhile.
go mod download all
# This will always take a while, CGO is slow.
go build -v -o rokon .
# This command doesn't exist. It's a placeholder for the future.
sudo make install
```

2. Run the app in development mode:

```bash
wgo go run -v .
```

Note: To truly test the development version of the app, you will need a Roku device.

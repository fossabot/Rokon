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

## Development Setup for Fedora Linux 40+

1. Clone the repository & install dependencies:

```bash
sudo dnf in -y git golang gtk4-devel @development-tools
git clone https://github.com/BrycensRanch/Rokon
cd Rokon/old
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

Note: To test the development version of the app, you will need a Roku device.

## Building the App

To build the app without testing on a Roku device, run the following command:

```bash
pnpm build
```

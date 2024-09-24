<<<<<<< HEAD
SHELL := /bin/bash
=======
#!/usr/bin/make -f
SHELL := $(shell which bash)
>>>>>>> d80b1e7 (chore: make the makefile use bash wherever it is)
# Define the default install directory
PREFIX ?= /usr/local
<<<<<<< HEAD
=======
VERSION ?= $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
DATE := $(shell date -u +%Y-%m-%d)
PACKAGED ?= false
PACKAGEFORMAT ?=
EXTRALDFLAGS :=
EXTRAGOFLAGS :=
BUILDTAGS :=

# Define target binary
TARGET = rokon
>>>>>>> d2f1c75 (build(debian): initial package)

# Define install directories
<<<<<<< HEAD
BINDIR = $(PREFIX)/bin
DATADIR = $(PREFIX)/share/rokon
APPLICATIONSDIR = $(PREFIX)/applications
ICONDIR = $(PREFIX)/icons/hicolor
METAINFODIR = $(PREFIX)/metainfo

# Define target binary
TARGET = rokon

=======
BINDIR = $(DESTDIR)$(PREFIX)/bin
DATADIR = $(DESTDIR)$(PREFIX)/share/rokon
DOCDIR = $(DESTDIR)$(PREFIX)/share/doc/rokon
APPLICATIONSDIR = $(DESTDIR)$(PREFIX)/share/applications
ICONDIR = $(DESTDIR)$(PREFIX)/share/icons/hicolor
METAINFODIR = $(DESTDIR)$(PREFIX)/share/metainfo
>>>>>>> 21b3760 (build(debian): standardize makefile)

.DEFAULT_GOAL := all
.PHONY: all
all: ## build pipeline
all: mod inst gen build spell lint test

.PHONY: precommit
precommit: ## validate the branch before commit
precommit: all vuln

.PHONY: ci
ci: ## CI build pipeline
ci: precommit diff

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## remove files created during build pipeline
	$(call print-target)
<<<<<<< HEAD
<<<<<<< HEAD
	rm -rf dist
	rm -rf pkg/
	rm *.pkg.tar.zst
	rm *.snap
	rm -f coverage.*
	rm -f '"$(shell go env GOCACHE)/../golangci-lint"'
=======
	rm -rf dist .flatpak flathub/.flatpak-builder flathub/repo AppDir pkg/ *.pkg.tar.zst *.snap coverage.* '"$(shell go env GOCACHE)/../golangci-lint"'
>>>>>>> d2f1c75 (build(debian): initial package)
=======
	rm -rf dist .flatpak flathub/.flatpak-builder flathub/repo AppDir *.AppImage *.rpm *.deb *.msi *.exe pkg/ *.pkg.tar.zst *.snap coverage.* '"$(shell go env GOCACHE)/../golangci-lint"'
>>>>>>> 6048fb1 (build: clean even more built binaries on clean command)
	go clean -i -cache -testcache -modcache -fuzzcache -x

<<<<<<< HEAD
=======
.PHONY: version
version: ## software version e.g 1.0.0
	@echo "Version: $(VERSION)"

.PHONY: appimage
appimage: ## build AppImage using appimage-builder
	$(call print-target)
	@echo "Building AppImage version: $(VERSION)"
	rm -rf AppDir
	$(MAKE) PACKAGED=true PACKAGEFORMAT=AppImage EXTRAGOFLAGS="-trimpath" EXTRALDFLAGS="-s -w" build
	$(MAKE) PREFIX=AppDir/usr BINDIR=AppDir install
	VERSION=$(VERSION) APPIMAGELAUNCHER_DISABLE=1 appimage-builder

.PHONY: fatimage
fatimage: ## build self contained AppImage that can run on older Linux systems while CI is on development branch
	$(call print-target)
	@echo "Building AppImage version: $(VERSION) (FAT)"
	rm -rf AppDir
	$(MAKE) PACKAGED=true PACKAGEFORMAT=AppImage EXTRAGOFLAGS="-trimpath" EXTRALDFLAGS="-s -w" build
	$(MAKE) PREFIX=AppDir/usr install
	VERSION=$(VERSION) APPIMAGELAUNCHER_DISABLE=1 appimagetool -s deploy ./AppDir/usr/share/applications/io.github.brycensranch.Rokon.desktop
	# My application follows the https://docs.fedoraproject.org/en-US/packaging-guidelines/AppData/ but this tool doesn't care lol
	mv AppDir/usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml AppDir/usr/share/metainfo/io.github.brycensranch.Rokon.appdata.xml
	cp ./AppDir/usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png ./AppDir
	VERSION=$(VERSION) APPIMAGELAUNCHER_DISABLE=1 mkappimage -u "gh-releases-zsync|BrycensRanch|Rokon|latest|Rokon-*x86_64.AppImage.zsync" ./AppDir




>>>>>>> 50639a3 (build(appimage): actually put build properties on right commands)
.PHONY: mod
mod: ## go mod tidy
	$(call print-target)
	go mod tidy
	cd tools && go mod tidy

.PHONY: inst
inst: ## go install tools
	$(call print-target)
	cd tools && go install $(shell cd tools && go list -e -f '{{ join .Imports " " }}' -tags=tools)

## Linux only, I have no idea how to do this on Windows
.PHONY: install
install:
	$(call print-target)
	@echo "Installing $(TARGET) to $(BINDIR)"
	mkdir -p $(BINDIR)
	install -Dpm 0755 $(TARGET) $(BINDIR)
<<<<<<< HEAD
	desktop-file-install --dir=$(PREFIX)/share/applications ./usr/share/applications/io.github.brycensranch.Rokon.desktop
	install -Dm644 ./usr/share/applications/io.github.brycensranch.Rokon.service $(PREFIX)/share/dbus-1/services/io.github.brycensranch.Rokon.service
	install -Dpm 0644 ./usr/share/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png $(PREFIX)/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png
	install -Dpm 0644 ./usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png $(PREFIX)/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png
	install -Dpm 0644 ./usr/share/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg $(PREFIX)/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg
	install -Dpm 0644 ./usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml $(PREFIX)/metainfo/io.github.brycensranch.Rokon.metainfo.xml
	update-desktop-database
=======
	install -Dpm 0644 ./usr/share/applications/io.github.brycensranch.Rokon.desktop $(APPLICATIONSDIR)/io.github.brycensranch.Rokon.desktop
	install -Dpm 0644 ./usr/share/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png $(ICONDIR)/48x48/apps/io.github.brycensranch.Rokon.png
	install -Dpm 0644 ./usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png $(ICONDIR)/256x256/apps/io.github.brycensranch.Rokon.png
	install -Dpm 0644 ./usr/share/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg $(ICONDIR)/scalable/apps/io.github.brycensranch.Rokon.svg
	install -Dpm 0644 ./usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml $(METAINFODIR)/io.github.brycensranch.Rokon.metainfo.xml
	install -Dpm 0644 ./LICENSE.md $(DESTDIR)$(PREFIX)/share/licenses/rokon/LICENSE.md
	install -Dpm 0644 ./README.md $(DESTDIR)$(PREFIX)/share/doc/rokon/README.md
	install -Dpm 0644 ./PRIVACY.md $(DESTDIR)$(PREFIX)/share/doc/rokon/PRIVACY.md
>>>>>>> 21b3760 (build(debian): standardize makefile)

.PHONY: uninstall
uninstall:
	$(call print-target)
	rm -f $(BINDIR)/$(TARGET)
<<<<<<< HEAD
	rm -f $(PREFIX)/share/dbus-1/services/io.github.brycensranch.Rokon.service
	rm -f $(PREFIX)/share/applications/io.github.brycensranch.Rokon.desktop
	rm -f $(PREFIX)/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png
	rm -f $(PREFIX)/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png
	rm -f $(PREFIX)/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg
	rm -f $(PREFIX)/metainfo/io.github.brycensranch.Rokon.metainfo.xml
	update-desktop-database
=======
	rm -f $(APPLICATIONSDIR)/io.github.brycensranch.Rokon.desktop
	rm -f $(ICONDIR)/48x48/apps/io.github.brycensranch.Rokon.png
	rm -f $(ICONDIR)/256x256/apps/io.github.brycensranch.Rokon.png
	rm -f $(ICONDIR)/scalable/apps/io.github.brycensranch.Rokon.svg
	rm -f $(METAINFODIR)/io.github.brycensranch.Rokon.metainfo.xml
	rm -rf $(DOCDIR)/rokon
	rm -rf $(DESTDIR)$(PREFIX)/share/licenses/rokon
>>>>>>> 21b3760 (build(debian): standardize makefile)

.PHONY: gen
gen: ## go generate
	$(call print-target)
	go generate ./...

.PHONY: build
build: ## go build -v -o rokon
	$(call print-target)
	go build -v -o $(TARGET)

.PHONY: spell
spell: ## misspell
	$(call print-target)
	misspell -error -locale=US -w **.md

.PHONY: lint
lint: ## golangci-lint
	$(call print-target)
	golangci-lint run --fix

.PHONY: vuln
vuln: ## govulncheck
	$(call print-target)
	govulncheck ./...

.PHONY: test
test: ## go test
	$(call print-target)
	go test -race -covermode=atomic -coverprofile=coverage.out -coverpkg=./... ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: diff
diff: ## git diff
	$(call print-target)
	git diff --exit-code
	RES=$$(git status --porcelain) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi


define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef

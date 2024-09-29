#!/usr/bin/make -f
SHELL := $(shell which bash)
# Define the default install directory
PREFIX ?= /usr/local
VERSION ?= $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
DATE := $(shell date -u +%Y-%m-%d)
PACKAGED ?= false
PACKAGEFORMAT ?=
NODOCUMENTATION ?= 0
EXTRALDFLAGS :=
EXTRAGOFLAGS :=
BUILDTAGS :=


ifneq ($(CFLAGS),)
    export CGO_CFLAGS := $(CFLAGS)
    $(info Using provided CFLAGS: $(CFLAGS))
endif

ifneq ($(CPPFLAGS),)
    export CGO_CPPFLAGS := $(CPPFLAGS)
    $(info Using provided CPPFLAGS: $(CPPFLAGS))
endif

ifneq ($(CXXFLAGS),)
    export CGO_CXXFLAGS := $(CXXFLAGS)
    $(info Using provided CXXFLAGS: $(CXXFLAGS))
endif

ifneq ($(LDFLAGS),)
    export CGO_LDFLAGS := $(LDFLAGS)
    $(info Using provided LDFLAGS: $(LDFLAGS))
endif

# Define target binary
TARGET = rokon

# Define install directories
BINDIR = $(DESTDIR)$(PREFIX)/bin
DATADIR = $(DESTDIR)$(PREFIX)/share/rokon
DOCDIR = $(DESTDIR)$(PREFIX)/share/doc/rokon
APPLICATIONSDIR = $(DESTDIR)$(PREFIX)/share/applications
ICONDIR = $(DESTDIR)$(PREFIX)/share/icons/hicolor
METAINFODIR = $(DESTDIR)$(PREFIX)/share/metainfo

.DEFAULT_GOAL := build
.PHONY: all
all: ## build pipeline
all: mod inst gen build spell lint test

.PHONY: precommit
precommit: ## validate the branch before commit
precommit: all vuln

.PHONY: ci
ci: ## CI checks pipeline
ci: precommit diff

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## remove files created during build pipeline
	$(call print-target)
	rm -rf dist .flatpak flathub/.flatpak-builder flathub/repo AppDir *.AppImage *.rpm *.deb *.msi *.exe pkg/ *.pkg.tar.zst *.snap coverage.* '"$(shell go env GOCACHE)/../golangci-lint"'
	go clean -i -cache -testcache -modcache -fuzzcache -x

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
	@echo "version $(VERSION)"
	mkdir -p $(BINDIR)
	mkdir -p $(DESTDIR)$(PREFIX)/share/doc/rokon
	install -Dpm 0755 $(TARGET) $(BINDIR) || true
	desktop-file-install --dir=$(DESTDIR)$(PREFIX)/share/applications ./usr/share/applications/io.github.brycensranch.Rokon.desktop
	install -Dm644 ./usr/share/applications/io.github.brycensranch.Rokon.service $(PREFIX)/share/dbus-1/services/io.github.brycensranch.Rokon.service || true
	install -Dpm 0644 ./usr/share/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png $(ICONDIR)/48x48/apps/io.github.brycensranch.Rokon.png
	install -Dpm 0644 ./usr/share/icons/hicolor/128x128/apps/io.github.brycensranch.Rokon.png $(ICONDIR)/48x48/apps/io.github.brycensranch.Rokon.png
	install -Dpm 0644 ./usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png $(ICONDIR)/256x256/apps/io.github.brycensranch.Rokon.png
	install -Dpm 0644 ./usr/share/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg $(ICONDIR)/scalable/apps/io.github.brycensranch.Rokon.svg
	install -Dpm 0644 ./usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml $(METAINFODIR)/io.github.brycensranch.Rokon.metainfo.xml
	install -Dpm 0644 ./LICENSE.md $(DESTDIR)$(PREFIX)/share/licenses/rokon/LICENSE.md
	# Check if NODOCUMENTATION is set to 1
	@if [ "$(NODOCUMENTATION)" != "1" ]; then \
		install -Dpm 0644 ./PRIVACY.md ./README.md $(DESTDIR)$(PREFIX)/share/doc/rokon; \
	else \
		echo "Skipping documentation installation. Please make sure you include PRIVACY notice."; \
	fi
	update-desktop-database || true

.PHONY: uninstall
uninstall:
	$(call print-target)
	@echo "Uninstalling version $(VERSION)"
	rm -f $(BINDIR)/$(TARGET)
	rm -f $(APPLICATIONSDIR)/io.github.brycensranch.Rokon.desktop
	rm -f $(ICONDIR)/48x48/apps/io.github.brycensranch.Rokon.png
	rm -f $(ICONDIR)/128x128/apps/io.github.brycensranch.Rokon.png
	rm -f $(ICONDIR)/256x256/apps/io.github.brycensranch.Rokon.png
	rm -f $(ICONDIR)/scalable/apps/io.github.brycensranch.Rokon.svg
	rm -f $(METAINFODIR)/io.github.brycensranch.Rokon.metainfo.xml
	rm -rf $(DOCDIR)/rokon
	rm -rf $(DESTDIR)$(PREFIX)/share/licenses/rokon

.PHONY: gen
gen: ## go generate
	$(call print-target)
	go generate ./...

.PHONY: build
build: ## go build -v -o rokon
	$(call print-target)
	@echo "Building version $(VERSION)"
	go build -v -ldflags="-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.packaged=$(PACKAGED) -X main.packageFormat=$(PACKAGEFORMAT) -X main.branch=$(BRANCH) $(EXTRALDFLAGS)" $(EXTRAGOFLAGS) -o $(TARGET) -tags "$(BUILDTAGS)" .

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

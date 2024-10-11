#!/usr/bin/make -f
SHELL := $(shell which bash)
# Define the default install directory
PREFIX ?= /usr/local
VERSION ?= $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git branch --show-current)
DATE := $(shell date -u +%Y-%m-%d)
PACKAGED ?= false
PACKAGEFORMAT ?=
NODOCUMENTATION ?= 0
EXTRALDFLAGS :=
EXTRAGOFLAGS :=
BUILDTAGS :=
UNAME_S := $(shell uname -s)



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
LICENSEDIR = $(DESTDIR)$(PREFIX)/share/licenses/rokon
APPLICATIONSDIR = $(DESTDIR)$(PREFIX)/share/applications
ICONDIR = $(DESTDIR)$(PREFIX)/share/icons/hicolor
METAINFODIR = $(DESTDIR)$(PREFIX)/share/metainfo
TARBALLDIR ?= ./tarball
ABS_TARBALLDIR := $(shell realpath $(TARBALLDIR))

LIBS_DIR ?= $(TARBALLDIR)/libs
TAR_NAME ?= rokon-$(shell uname)-$(VERSION)-$(shell uname -m).tar.gz


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
	rm -rf dist .flatpak io.github.brycensranch.Rokon.desktop tarball io.github.brycensranch.Rokon.metainfo.xml macos/rokon .flatpak-builder flathub/.flatpak-builder flathub/repo flathub/export macos/share flathub/*.flatpak AppDir *.AppImage *.rpm *.pdf *.rtf windows/*.rtf *.deb *.msi *.exe pkg/ *.pkg.tar.zst .ptmp* *.tar* *.snap *.zsync rokon debian/tmp debian/rokon debian/.debhelper debian/rokon-dbg coverage.* '"$(shell go env GOCACHE)/../golangci-lint"'
	# go clean -i -cache -testcache -modcache -fuzzcache -x

.PHONY: nuke
nuke: ## completely clean the repository of artifacts and clear cache
	$(call print-target)
	$(MAKE) clean
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
	rm AppDir/usr/lib64/libLLVM* || true
	# My application follows the https://docs.fedoraproject.org/en-US/packaging-guidelines/AppData/ but this tool doesn't care lol
	mv AppDir/usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml AppDir/usr/share/metainfo/io.github.brycensranch.Rokon.appdata.xml
	cp ./AppDir/usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png ./AppDir
	VERSION=$(VERSION) APPIMAGELAUNCHER_DISABLE=1 mkappimage -u "gh-releases-zsync|BrycensRanch|Rokon|latest|Rokon-*$(shell uname -m).AppImage.zsync" ./AppDir

.PHONY: tarball
tarball: ## build self contained Tarball that auto updates
	$(call print-target)
	@echo "Building Rokon Tarball version: $(VERSION)"
	rm -v -rf $(TARBALLDIR) || sudo rm -rf $(TARBALLDIR)
	mkdir -p $(TARBALLDIR)
	mkdir -p $(LIBS_DIR)
	$(MAKE) PACKAGED=true PACKAGEFORMAT=portable EXTRAGOFLAGS="-trimpath -buildmode=pie" EXTRALDFLAGS="-s -w -linkmode=external" build
	$(MAKE) PREFIX=$(TARBALLDIR) BINDIR=$(TARBALLDIR) APPLICATIONSDIR=$(TARBALLDIR) install
	cp -v ./windows/portable.txt $(TARBALLDIR)
	ldd -d -r $(TARGET) | awk '{print $$3}' | grep -v 'not found' | while read -r dep; do \
		cp -L --no-preserve=mode --debug "$$dep" $(LIBS_DIR); \
	done
	@cp -L --no-preserve=mode --debug $$(ldd ./rokon | grep 'ld-linux' | awk '{print $$1}') $(TARBALLDIR)/libs/
	@chmod +x $(LIBS_DIR)/*.so*
	@strip -p -v $(LIBS_DIR)/*.so*
	patchelf --force-rpath --set-rpath ./libs $(TARBALLDIR)/$(TARGET)
	echo '#!/bin/sh' > $(TARBALLDIR)/rokon.sh; \
	echo 'export LD_LIBRARY_PATH="./libs:$${LD_LIBRARY_PATH}"' >> $(TARBALLDIR)/rokon.sh; \
	echo 'export LD_PRELOAD="./libs/libc.so.6"' >> $(TARBALLDIR)/rokon.sh; \
	echo 'export XKB_DEFAULT_INCLUDE_PATH=./share/X11/xkb' >> $(TARBALLDIR)/rokon.sh; \
	echo 'export XKB_CONFIG_ROOT=./share/X11/xkb' >> $(TARBALLDIR)/rokon.sh; \
	echo 'exec ./libs/$(shell /bin/ls $(LIBS_DIR) | grep "ld-linux") ./$(TARGET) "$$@"' >> $(TARBALLDIR)/rokon.sh; \
	chmod +x $(TARBALLDIR)/rokon.sh
	cd /usr && cp -r --parents -L -v --no-preserve=mode -r share/glib-2.0/schemas/gschemas.compiled share/X11 share/gtk-4.0 share/icons/Adwaita share/icons/AdwaitaLegacy lib/gdk-pixbuf-2.0 $(ABS_TARBALLDIR)
	sed -i 's/rokon/\.\/rokon.sh/g' $(TARBALLDIR)/io.github.brycensranch.Rokon.desktop
	cd $(TARBALLDIR) && ./rokon.sh --version # sanity check
	tar -czf $(TAR_NAME) $(TARBALLDIR)
	@if command -v zsyncmake >/dev/null 2>&1; then \
		zsyncmake $(TAR_NAME) -u "gh-releases-zsync|BrycensRanch|Rokon|latest|Rokon-$(shell uname)-*-$(shell uname -m).tar.gz.zsync"; \
	else \
		@echo "zsyncmake not found. Please install it to generate the zsync file."; \
	fi
	rm $(TARGET)
	@echo "Tarball created: $(TAR_NAME)"

.PHONY: dev
dev: ## go run -v .
	$(call print-target)
	@echo "Starting development server for Rokon: $(VERSION)"
	go run -v .

.PHONY: run
run: ## go run -v .
run: dev

.PHONY: mod
mod: ## go mod tidy
	$(call print-target)
	go mod tidy
	cd tools && go mod tidy

.PHONY: inst
inst: ## go install tools
	$(call print-target)
	cd tools && go get $(shell cd tools && go list -e -f '{{ join .Imports " " }}' -tags=tools)


.PHONY: install
install: ## installs Rokon into $PATH and places desktop files
	$(call print-target)
	@echo "Installing $(TARGET) to $(BINDIR)"
	@echo "Version: $(VERSION)"
	@echo "Creating necessary directories..."
	mkdir -p $(BINDIR)
	mkdir -p $(DESTDIR)$(PREFIX)/share/doc/rokon
	mkdir -p $(APPLICATIONSDIR)
	mkdir -p $(LICENSEDIR)
	mkdir -p $(ICONDIR)/48x48/apps
	mkdir -p $(ICONDIR)/256x256/apps
	mkdir -p $(ICONDIR)/scalable/apps
	mkdir -p $(METAINFODIR)
	@echo "Detected OS: $(UNAME_S)"

	# Install the target
	@echo "Installing target..."
	if [ "$(UNAME_S)" = "Darwin" ]; then \
		echo "Installing $(TARGET) to $(BINDIR) as macOS binary..."; \
		install -m 0755 $(TARGET) $(BINDIR) || true; \
	else \
		echo "Installing $(TARGET) to $(BINDIR) as Linux binary..."; \
		install -Dpm 0755 $(TARGET) $(BINDIR) || true; \
	fi

	# Install desktop files and icons
	@echo "Installing desktop files and icons..."
	if [ "$(UNAME_S)" = "Darwin" ]; then \
		echo "Installing desktop file to $(APPLICATIONSDIR)/io.github.brycensranch.Rokon.desktop"; \
		install -m 0644 ./usr/share/applications/io.github.brycensranch.Rokon.desktop $(APPLICATIONSDIR)/io.github.brycensranch.Rokon.desktop; \
		install -m 0644 ./usr/share/dbus-1/services/io.github.brycensranch.Rokon.service $(DESTDIR)$(PREFIX)/share/dbus-1/services/io.github.brycensranch.Rokon.service; \
		echo "Installing icon to $(ICONDIR)/48x48/apps/io.github.brycensranch.Rokon.png"; \
		install -m 0644 ./usr/share/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png $(ICONDIR)/48x48/apps/io.github.brycensranch.Rokon.png; \
		echo "Installing icon to $(ICONDIR)/128x128/apps/io.github.brycensranch.Rokon.png"; \
		install -m 0644 ./usr/share/icons/hicolor/128x128/apps/io.github.brycensranch.Rokon.png $(ICONDIR)/128x128/apps/io.github.brycensranch.Rokon.png; \
		echo "Installing icon to $(ICONDIR)/256x256/apps/io.github.brycensranch.Rokon.png"; \
		install -m 0644 ./usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png $(ICONDIR)/256x256/apps/io.github.brycensranch.Rokon.png; \
		echo "Installing SVG icon to $(ICONDIR)/scalable/apps/io.github.brycensranch.Rokon.svg"; \
		install -m 0644 ./usr/share/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg $(ICONDIR)/scalable/apps/io.github.brycensranch.Rokon.svg; \
		echo "Installing metainfo file to $(METAINFODIR)/io.github.brycensranch.Rokon.metainfo.xml"; \
		install -m 0644 ./usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml $(METAINFODIR)/io.github.brycensranch.Rokon.metainfo.xml; \
		echo "Installing license to $(LICENSEDIR)/LICENSE.md"; \
		install -m 0644 ./LICENSE.md $(LICENSEDIR)/LICENSE.md; \
	else \
		echo "Installing desktop file to $(APPLICATIONSDIR)/io.github.brycensranch.Rokon.desktop"; \
		install -Dpm 0644 ./usr/share/applications/io.github.brycensranch.Rokon.desktop $(APPLICATIONSDIR)/io.github.brycensranch.Rokon.desktop; \
		install -Dpm 0644 ./usr/share/dbus-1/services/io.github.brycensranch.Rokon.service $(DESTDIR)$(PREFIX)/share/dbus-1/services/io.github.brycensranch.Rokon.service; \
		echo "Installing icon to $(ICONDIR)/48x48/apps/io.github.brycensranch.Rokon.png"; \
		install -Dpm 0644 ./usr/share/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png $(ICONDIR)/48x48/apps/io.github.brycensranch.Rokon.png; \
		echo "Installing icon to $(ICONDIR)/128x128/apps/io.github.brycensranch.Rokon.png"; \
		install -Dpm 0644 ./usr/share/icons/hicolor/128x128/apps/io.github.brycensranch.Rokon.png $(ICONDIR)/128x128/apps/io.github.brycensranch.Rokon.png; \
		echo "Installing icon to $(ICONDIR)/256x256/apps/io.github.brycensranch.Rokon.png"; \
		install -Dpm 0644 ./usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png $(ICONDIR)/256x256/apps/io.github.brycensranch.Rokon.png; \
		echo "Installing SVG icon to $(ICONDIR)/scalable/apps/io.github.brycensranch.Rokon.svg"; \
		install -Dpm 0644 ./usr/share/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg $(ICONDIR)/scalable/apps/io.github.brycensranch.Rokon.svg; \
		echo "Installing metainfo file to $(METAINFODIR)/io.github.brycensranch.Rokon.metainfo.xml"; \
		install -Dpm 0644 ./usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml $(METAINFODIR)/io.github.brycensranch.Rokon.metainfo.xml; \
		echo "Installing license to $(LICENSEDIR)/LICENSE.md"; \
		install -Dpm 0644 ./LICENSE.md $(LICENSEDIR)/LICENSE.md; \
	fi

	# Check if NODOCUMENTATION is set to 1
	@if [ "$(NODOCUMENTATION)" != "1" ]; then \
		echo "Installing documentation..."; \
		if [ "$(UNAME_S)" = "Darwin" ]; then \
			install -m 0644 ./PRIVACY.md ./README.md $(DESTDIR)$(PREFIX)/share/doc/rokon; \
		else \
			install -Dpm 0644 ./PRIVACY.md ./README.md $(DESTDIR)$(PREFIX)/share/doc/rokon; \
		fi; \
	else \
		echo "Skipping documentation installation. Please make sure you include PRIVACY notice."; \
	fi


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
	rm -f $(DESTDIR)$(PREFIX)/share/dbus-1/services/io.github.brycensranch.Rokon.service
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

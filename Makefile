SHELL := /bin/bash

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
	rm -rf dist
	rm -f coverage.*
	rm -f '"$(shell go env GOCACHE)/../golangci-lint"'
	go clean -i -cache -testcache -modcache -fuzzcache -x

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
	install -Dpm 0755 ./rokon /usr/local/bin/rokon
	install -Dpm 0644 ./usr/share/applications/io.github.brycensranch.Rokon.desktop /usr/share/applications/io.github.brycensranch.Rokon.desktop
	install -Dpm 0644 ./usr/share/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png /usr/share/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png
	install -Dpm 0644 ./usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png /usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png
	install -Dpm 0644 ./usr/share/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg /usr/share/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg
	install -Dpm 0644 ./usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml /usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml

.PHONY: uninstall
uninstall:
	$(call print-target)
	rm -f /usr/local/bin/rokon
	rm -f /usr/share/applications/io.github.brycensranch.Rokon.desktop
	rm -f /usr/share/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png
	rm -f /usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png
	rm -f /usr/share/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg
	rm -f /usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml

.PHONY: gen
gen: ## go generate
	$(call print-target)
	go generate ./...

.PHONY: build
build: ## goreleaser build
	$(call print-target)
	goreleaser build --single-target --snapshot

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

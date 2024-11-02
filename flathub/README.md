# Rokon Flatpak

This directory is meant to be mirrored to https://github.com/flathub/io.github.brycensranch.Rokon

## Test changes

```sh
flatpak run org.flatpak.Builder --force-clean --sandbox --user --install --install-deps-from=flathub --ccache --mirror-screenshots-url=https://dl.flathub.org/media/ --repo=repo builddir io.github.brycensranch.Rokon.yml
```

## Run it!

Rokon's stable flatpak was just built and automatically installed to the master branch of your local flatpak repository.
So, let's test it!

`flatpak run --branch=master io.github.brycensranch.Rokon`

## Wondering how to create a flatpak bundle?

```sh
# After building with the command provided above...
flatpak build-export export builddir
flatpak build-bundle export io.github.brycensranch.Rokon.flatpak io.github.brycensranch.Rokon master
```

## Important for developers

After you're done with your changes, you should go into the root of this repository and run `make clean`. This will get rid of the builddir, flatpak bundles, repos, etc. This is needed because on my machine, all these nested directories makes Golang's language server very very unhappy. At least on VSCode. :3

## Note

This repository is automatically managed by GitHub Actions. This includes: Go dependencies, commit SHAs.

The "beta" flatpak with Rokon from git will never have a different ID. io.github.brycensranch.Rokon-beta is just a placeholder. To build the beta flatpak, simply replace any mentions of io.github.brycensranch.Rokon.yml to io.github.brycensranch.Rokon-beta.yml.

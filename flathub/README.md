# Rokon Flatpak

This directory is meant to be mirrored to https://github.com/flathub/io.github.brycensranch.Rokon

## Test changes

```sh
flatpak run org.flatpak.Builder --force-clean --sandbox --user --install --install-deps-from=flathub --ccache --mirror-screenshots-url=https://dl.flathub.org/media/ --repo=repo builddir io.github.brycensranch.Rokon
```

## Note

This repository is automatically managed by GitHub Actions. This includes: Go dependencies, commit SHAs.

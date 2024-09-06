# Come on, LOCK IN!!
# Maintainer: Brycen Granville <brycengranville@outlook.com>

pkgname=rokon
pkgver=1.0.0+7d4e7cb
pkgrel=1
epoch=0
maintainer="Brycen Granville <brycengranville@outlook.com>"
pkgdesc="A roku remote for your desktop"
arch=('x86_64' 'aarch64' 'armv7h')
url="https://github.com/BrycensRanch/Rokon"
license=('AGPL3-or-later')
depends=('gtk4')
makedepends=('gtk4' 'git' 'gcc' 'go' 'cmake')
source=("git+https://github.com/BrycensRanch/Rokon.git#branch=master")
sha256sums=('SKIP')

pkgver() {
    cd Rokon
    # Get the latest commit hash
    local commit_hash=$(git rev-parse --short HEAD)

    # Extract the base version before the '+'
    local base_ver=$(echo "$pkgver" | cut -d'+' -f1)

    # Return the base version with the updated commit hash
    echo "${base_ver}+${commit_hash}"
}

changelog() {
    cd Rokon
    git log --pretty=format:'%ad %h %s' --date=short
}

build() {
    cd Rokon
    go mod download all
    go build -v -ldflags="-X main.version=$pkgver -X main.commit=$(git rev-parse --short HEAD) -X main.packaged=true -X main.packageFormat=arch -X main.branch=$(git rev-parse --abbrev-ref HEAD) -X main.date=$(date -u +%Y-%m-%d)" -o $pkgname
}

package() {
    cd Rokon
    install -Dm755 $pkgname $pkgdir/usr/bin/$pkgname
    install -Dm644 LICENSE.md $pkgdir/usr/share/licenses/$pkgname/LICENSE.md
    install -Dm644 README.md $pkgdir/usr/share/doc/$pkgname/README.md
    install -Dpm 0644 ./usr/share/applications/io.github.brycensranch.Rokon.desktop $pkgdir/usr/share/applications/io.github.brycensranch.Rokon.desktop
    install -Dpm 0644 ./usr/share/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png $pkgdir/usr/share/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png
    install -Dpm 0644 ./usr/share/icons/hicolor/128x128/apps/io.github.brycensranch.Rokon.png $pkgdir/usr/share/icons/hicolor/128x128/apps/io.github.brycensranch.Rokon.png
    install -Dpm 0644 ./usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png $pkgdir/usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png
    install -Dpm 0644 ./usr/share/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg $pkgdir/usr/share/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg
    install -Dpm 0644 ./usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml $pkgdir/usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml
}


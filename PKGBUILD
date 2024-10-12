# Maintainer: Brycen Granville <brycengranville@outlook.com>

pkgname=rokon
pkgver=1.0.0+7d4e7cb
pkgrel=1
epoch=0
maintainer="Brycen Granville <brycengranville@outlook.com>"
pkgdesc="A roku remote for your desktop"
arch=('x86_64' 'aarch64' 'armv7h')
url="https://github.com/BrycensRanch/Rokon"
license=('AGPL-3.0-or-later')
depends=('gtk4')
makedepends=('gtk4' 'git' 'gcc' 'go' 'make')
source=("git+https://github.com/BrycensRanch/Rokon.git#branch=master")
sha256sums=('SKIP')

pkgver() {
    cd Rokon
    # Get the latest commit hash
    local commit_hash=$(git rev-parse --short HEAD)

    local base_ver=$(cat VERSION)

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
    make BRANCH=$(git rev-parse --abbrev-ref HEAD)  TARGET=$pkgname PACKAGED=true PACKAGEFORMAT=arch EXTRAGOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw" EXTRALDFLAGS="-compressdwarf=false -linkmode=external" build
}

package() {
    cd Rokon
    make TARGET=$pkgname PREFIX=$pkgdir/usr install
}


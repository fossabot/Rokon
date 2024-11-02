# Contributor: Brycen G <brycengranville@outlook.com>
# Maintainer: Brycen G <brycengranville@outlook.com>
pkgname=rokon
pkgver=1.0.0
pkgrel=0
pkgdesc="A roku remote for your desktop"
url="https://github.com/BrycensRanch/Rokon"
arch="all"
license="AGPL-3.0-or-later"
makedepends="
    gtk4.0-dev
    gobject-introspection-dev
    go
    build-base
    unzip
	"
subpackages="
	$pkgname-doc
	"
source="https://nightly.link/BrycensRanch/Rokon/actions/runs/11642705367/rokon-vendored-source.zip"


build() {
    make TARGET=$pkgname PACKAGED=true PACKAGEFORMAT=apk EXTRAGOFLAGS="-buildmode=pie -trimpath -mod=vendor -modcacherw" EXTRALDFLAGS="-compressdwarf=false -linkmode=external" build
}

check() {
    ./rokon --version
}

package() {
    cd src
    make DESTDIR="$pkgdir" PREFIX="/usr" install
}

sha512sums="
f225bb9e8fd11e97cfc0adb4b952dcc8c626fd01aa28e1d57f30322a777026acad4bc6f0319ea0ed228b1409348c829ad15581a224c05e97c5ac66c414f7a88e rokon-vendored-source.zip
"

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
    bash
    unzip
	"
subpackages="
	$pkgname-doc
	"
source="https://nightly.link/BrycensRanch/Rokon/actions/runs/11487693784/rokon-vendored-source.zip"


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
ffa45098ef49b9f6439575dc34cf195d3d3de9be2766aba4da4a19062ab835a16f64027d8002cb072afc1853d75657a2ac5d2e81beb5e1bacad0ab1607d232ce rokon-vendored-source.zip
"

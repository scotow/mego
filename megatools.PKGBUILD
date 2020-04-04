# Maintainer: Benjamin Lopez <contact@scotow.com>

pkgname=megatools
pkgver=1.11.0
pkgrel=20200404
pkgdesc="Command line client application for Mega.nz"
arch=('x86_64')
url="http://megatools.megous.com"
license=('GPL')
depends=('curl' 'glib2' 'openssl')
makedepends=('asciidoc')
_archive="megatools-${pkgver}-git-${pkgrel}-linux-x86_64"
source=("https://megatools.megous.com/builds/experimental/${_archive}.tar.gz")
options=(!libtool)
sha256sums=('e2795b6126ff9401a830dbae0006dbb77c48fd5dab788c8f1ffd1efc42db2f54')

package() {
	cd "${_archive}"

	# Bin
	install -Dm755 "${pkgname}" "${pkgdir}/usr/bin/${pkgname}"

	# License
	install -Dm644 "LICENSE" "${pkgdir}/usr/share/licenses/${pkgname}/LICENSE"
}

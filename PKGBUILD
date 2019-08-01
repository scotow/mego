# Maintainer: Benjamin Lopez <contact@scotow.com>

pkgname=mego
pkgver=1.0.3
pkgrel=1
pkgdesc="A simple megadl wrapper with auto-retry and download list"
arch=('x86_64')
url="https://github.com/scotow/${pkgname}"
license=('MIT')
depends=('megatools')
makedepends=('go' 'git')
source=("${pkgname}-${pkgver}.tar.gz::https://github.com/scotow/${pkgname}/archive/${pkgver}.tar.gz")
sha256sums=('b9b8435fb7dbc1dc4a6e778d35f2be492e9872bfc6752538c5a55ed15b8d81d3')

prepare(){
  mkdir -p src/github.com/scotow
  ln -rTsf "${pkgname}-${pkgver}" "src/github.com/scotow/${pkgname}"
}

build(){
  export GOPATH="${srcdir}"
  cd "src/github.com/scotow/${pkgname}"
  go install \
	-gcflags "all=-trimpath=${GOPATH}/src" \
	-asmflags "all=-trimpath=${GOPATH}/src" \
	-ldflags "-extldflags ${LDFLAGS}" \
	./...
}

package(){
  install -Dm755 "bin/${pkgname}" "${pkgdir}/usr/bin/${pkgname}"

  cd "${pkgname}-${pkgver}"
  install -Dm644 LICENSE "${pkgdir}/usr/share/licenses/${pkgname}/LICENSE"
}

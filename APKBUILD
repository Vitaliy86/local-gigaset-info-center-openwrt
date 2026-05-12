# Maintainer: Vitaliy86 <vitaliy86@github.com>
pkgname="gigaset-info-center"
pkgver="2.0"
pkgrel=0
arch="aarch64 armv7 x86_64 mips_sf mipsel_sf"
pkgdesc="Replacement weather service for Gigaset IP handsets"
url="https://github.com/Vitaliy86/local-gigaset-info-center-openwrt"
license="AGPL-3.0-or-later"
depends=""
options="!strip !check !fhs"
# Binary is built by CI and placed in $startdir before abuild runs.
source=""

build() { return 0; }

package() {
    # Binary (arch-specific, pre-compiled by Go toolchain in CI)
    install -D -m 0755 "$startdir/gigaset-info-center" \
        "$pkgdir/usr/bin/gigaset-info-center"

    # Example config
    install -D -m 0644 "$startdir/gigaset-info-center.conf.example" \
        "$pkgdir/etc/gigaset-info-center.conf.example"

    # OpenWrt procd init script
    install -D -m 0755 "$startdir/gigaset-info-center.init" \
        "$pkgdir/etc/init.d/gigaset-info-center"
}

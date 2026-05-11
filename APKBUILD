# Maintainer: Vitaliy86 <vitaliy86@github.com>
pkgname="gigaset-info-center"
pkgver="1.7"
pkgrel=0
arch="noarch"
license="AGPL-3.0-or-later"
url="https://github.com/Vitaliy86/local-gigaset-info-center-openwrt"
maintainer="Vitaliy86 <vitaliy86@github.com>"
pkgdesc="Replacement weather service for Gigaset IP handsets"
depends="php8 lighttpd php8-mod-curl php8-mod-gd lighttpd-mod-fastcgi"
options="!strip !check !fhs"
source=""

build() { return 0; }

package() {
    # Веб-файлы — /usr/share вместо /srv (Alpine запрещает /srv)
    mkdir -p "$pkgdir/usr/share/gigaset-info-center/icons"
    cp "$startdir/index.php"   "$pkgdir/usr/share/gigaset-info-center/"
    cp "$startdir/weather.php" "$pkgdir/usr/share/gigaset-info-center/"
    cp "$startdir/proxy.php"   "$pkgdir/usr/share/gigaset-info-center/"
    cp "$startdir"/icons/*.png "$pkgdir/usr/share/gigaset-info-center/icons/"

    mkdir -p "$pkgdir/etc/lighttpd"
    cp "$startdir/etc/lighttpd/gigaset-info-center.conf" \
       "$pkgdir/etc/lighttpd/"
    cp "$startdir/etc/gigaset-env.example" "$pkgdir/etc/"

    mkdir -p "$pkgdir/etc/init.d"
    install -m 0755 "$startdir/gigaset-info-center.init" \
        "$pkgdir/etc/init.d/gigaset-info-center"
}

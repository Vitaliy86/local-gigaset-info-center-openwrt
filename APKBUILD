# This file is the build script for Alpine Linux packages.
# It follows the Alpine Packager guidelines.
# For OpenWRT 25.12.3+, this will produce a valid .apk package.

pkgname="gigaset-info-center"
_pkgver="1.7"
pkgver="${_pkgver}-r0"
epoch=0
arch="noarch"
license="AGPL-3.0-or-later"
url="https://github.com/Vitaliy86/local-gigaset-info-center-openwrt"
maintainer="Vitaliy86"
description="Replacement weather service for Gigaset IP handsets"
depends="php lighttpd php-curl php-gd"

noarch="noarch"
options="!strip"

build() {
    return 0
}

package() {
    local dest="$pkgdir"
    
    # Install web application files to /srv/gigaset-info-center/
    mkdir -p "$dest/srv/gigaset-info-center/icons"
    cp index.php "$dest/srv/gigaset-info-center/"
    cp weather.php "$dest/srv/gigaset-info-center/"
    cp proxy.php "$dest/srv/gigaset-info-center/"
    cp icons/*.png "$dest/srv/gigaset-info-center/icons/"
    
    # Install configuration files to /etc/
    mkdir -p "$dest/etc/lighttpd"
    cp etc/lighttpd/gigaset-info-center.conf "$dest/etc/"
    cp etc/gigaset-env.example "$dest/etc/"
    
    # Install init script to /etc/init.d/
    mkdir -p "$dest/etc/init.d"
    cp gigaset-info-center.init "$dest/etc/init.d/gigaset-info-center"
    chmod 0755 "$dest/etc/init.d/gigaset-info-center"
    
    # Install documentation
    mkdir -p "$dest/usr/share/doc/$pkgname"
    cp LICENSE README.md "$dest/usr/share/doc/$pkgname/"
}

md5sums="skip"

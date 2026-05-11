# Makefile for OpenWrt package: gigaset-info-center
# Based on https://openwrt.org/ru/doc/devel/packages
#
# Usage:
#   make package/gigaset-info-center/compile V=s
#   make package/gigaset-info-center/clean
#   make package/gigaset-info-center/install

include $(TOPDIR)/rules.mk

PKG_NAME:=gigaset-info-center
PKG_VERSION:=1.7
PKG_RELEASE:=1

PKG_LICENSE:=AGPL-3.0-or-later
PKG_LICENSE_FILES:=LICENSE
PKG_MAINTAINER:=Vitaliy86 <vitaliy86@github.com>

include $(INCLUDE_DIR)/package.mk

define Package/gigaset-info-center
  SECTION:=net
  CATEGORY:=Network
  TITLE:=Replacement weather service for Gigaset IP handsets
  URL:=https://github.com/Vitaliy86/local-gigaset-info-center-openwrt
  DEPENDS:=+php8 +lighttpd +php8-mod-curl +php8-mod-gd
endef

define Package/gigaset-info-center/description
  Replacement weather service for Gigaset IP handsets.
  Provides XHTML-based weather information display for Gigaset
  SIP handsets (SL55, SX503, etc.) using OpenWeatherMap API.
endef

define Package/gigaset-info-center/conffiles
/etc/lighttpd/gigaset-info-center.conf
/etc/gigaset-env
endef

define Build/Configure
endef

define Build/Compile
endef

define Package/gigaset-info-center/install
	$(INSTALL_DIR) $(1)/srv/gigaset-info-center
	$(CP) index.php $(1)/srv/gigaset-info-center/
	$(CP) weather.php $(1)/srv/gigaset-info-center/
	$(CP) proxy.php $(1)/srv/gigaset-info-center/
	$(INSTALL_DIR) $(1)/srv/gigaset-info-center/icons
	$(CP) icons/*.png $(1)/srv/gigaset-info-center/icons/
	$(INSTALL_DIR) $(1)/etc/lighttpd
	$(CP) etc/lighttpd/gigaset-info-center.conf $(1)/etc/lighttpd/
	$(INSTALL_DIR) $(1)/etc
	$(CP) etc/gigaset-env.example $(1)/etc/gigaset-env.example
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) gigaset-info-center.init $(1)/etc/init.d/gigaset-info-center
	$(INSTALL_DIR) $(1)/usr/share/doc/gigaset-info-center
	$(CP) LICENSE README.md $(1)/usr/share/doc/gigaset-info-center/
endef

$(eval $(call BuildPackage,gigaset-info-center))

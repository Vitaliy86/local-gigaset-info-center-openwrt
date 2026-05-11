# Makefile for OpenWrt package: gigaset-info-center
# Based on https://openwrt.org/ru/doc/devel/packages

include $(TOPDIR)/rules.mk
include $(INCLUDE_DIR)/package.mk

PKG_NAME:=gigaset-info-center
PKG_VERSION:=1.7
PKG_RELEASE:=1

PKG_LICENSE:=AGPL-3.0-or-later
PKG_LICENSE_FILES:=LICENSE
PKG_MAINTAINER:=Vitaliy86 <vitaliy86@github.com>

# Define package metadata
define Package/gigaset-info-center
  SECTION:=net
  CATEGORY:=Network
  TITLE:=Replacement weather service for Gigaset IP handsets
  URL:=https://github.com/Vitaliy86/local-gigaset-info-center-openwrt
  DEPENDS:=+php8 +lighttpd +lighttpd-mod-fastcgi +php8-mod-curl +php8-mod-gd
endef

define Package/gigaset-info-center/description
  Replacement weather service for Gigaset IP handsets.
  Provides XHTML-based weather information display for Gigaset
  SIP handsets (SL55, SX503, etc.) using OpenWeatherMap API.
endef

define Package/gigaset-info-center/conffiles
/etc/lighttpd/gigaset-info-center.conf
/etc/gigaset-env.example
endef

define Build/Configure
endef

define Package/gigaset-info-center/install
	$(INSTALL_DIR) $(1)/srv/gigaset-info-center
	$(CP) $(PKG_BUILD_DIR)/index.php $(1)/srv/gigaset-info-center/
	$(CP) $(PKG_BUILD_DIR)/weather.php $(1)/srv/gigaset-info-center/
	$(CP) $(PKG_BUILD_DIR)/proxy.php $(1)/srv/gigaset-info-center/
	$(INSTALL_DIR) $(1)/srv/gigaset-info-center/icons
	$(CP) $(PKG_BUILD_DIR)/icons/*.png $(1)/srv/gigaset-info-center/icons/
	$(INSTALL_DIR) $(1)/etc/lighttpd
	$(CP) $(PKG_BUILD_DIR)/etc/lighttpd/gigaset-info-center.conf $(1)/etc/lighttpd/
	$(INSTALL_DIR) $(1)/etc
	$(CP) $(PKG_BUILD_DIR)/etc/gigaset-env.example $(1)/etc/gigaset-env.example
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/gigaset-info-center.init $(1)/etc/init.d/gigaset-info-center
	$(INSTALL_DIR) $(1)/usr/share/doc/gigaset-info-center
	$(CP) $(PKG_BUILD_DIR)/LICENSE $(PKG_BUILD_DIR)/README.md $(1)/usr/share/doc/gigaset-info-center/
endef

$(eval $(call BuildPackage,gigaset-info-center))

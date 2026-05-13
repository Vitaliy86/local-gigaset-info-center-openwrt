include $(TOPDIR)/rules.mk
include $(INCLUDE_DIR)/package.mk
include $(INCLUDE_DIR)/nmk.mk

PKG_NAME:=gigaset-info-center
PKG_VERSION:=2.0
PKG_RELEASE:=1

PKG_MAINTAINER:=Vitaliy86 <vitaliy86@github.com>
PKG_LICENSE:=AGPL-3.0-or-later
PKG_LICENSE_FILES:=LICENSE

PKG_SOURCE_PROTO:=git
PKG_SOURCE_URL:=https://github.com/Vitaliy86/local-gigaset-info-center-openwrt.git
PKG_SOURCE_VERSION:=HEAD
PKG_MIRROR_HASH:=skip

PKG_BUILD_DIR:=$(BUILD_DIR)/$(PKG_NAME)

include $(INCLUDE_DIR)/package.mk

define Package/gigaset-info-center
  SECTION:=utils
  CATEGORY:=Utilities
  TITLE:=Gigaset Info Center
  URL:=https://github.com/Vitaliy86/local-gigaset-info-center-openwrt
  DEPENDS:=+libc
endef

define Package/gigaset-info-center/description
  Replacement weather service for Gigaset IP handsets
endef

define Build/Prepare
	$(CP) $(PKG_BUILD_DIR)/* $(PKG_DIR)/
endef

define Build/Compile
endef

define Package/gigaset-info-center/install
	$(INSTALL_DIR) $(1)/usr/bin
	$(INSTALL_BIN) $(PKG_DIR)/gigaset-info-center $(1)/usr/bin/
	$(INSTALL_DIR) $(1)/etc
	$(INSTALL_CONF) $(PKG_DIR)/gigaset-info-center.conf.example $(1)/etc/gigaset-info-center.conf.example
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) $(PKG_DIR)/gigaset-info-center.init $(1)/etc/init.d/gigaset-info-center
endef

$(eval $(call BuildPackage,gigaset-info-center))

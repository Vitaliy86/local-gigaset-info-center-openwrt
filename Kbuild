include $(TOPDIR)/rules.mk

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
include $(INCLUDE_DIR)/cmake.mk

define Package/gigaset-info-center
  SECTION:=utils
  CATEGORY:=Utilities
  TITLE:=Gigaset Info Center
  URL:=https://github.com/Vitaliy86/local-gigaset-info-center-openwrt
endef

define Package/gigaset-info-center/description
  Replacement weather service for Gigaset IP handsets
endef

define Build/Configure
	$(call Build/Configure/CMake,$(PKG_BUILD_DIR))
endef

define Build/Compile
	$(call Build/Compile/CMake,$(PKG_BUILD_DIR),-DCMAKE_FIND_ROOT_PATH_MODE_PACKAGE=ONLY)
endef

define Package/gigaset-info-center/install
	$(INSTALL_DIR)$(1)/usr/bin
	$(INSTALL_BIN)$(PKG_BUILD_DIR)/gigaset-info-center$(1)/usr/bin/
endef

$(eval $(call BuildPackage,gigaset-info-center))

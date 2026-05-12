include $(TOPDIR)/rules.mk

PKG_NAME:=gigaset-info-center
PKG_VERSION:=2.0
PKG_RELEASE:=1

include $(INCLUDE_DIR)/package.mk

define Package/$(PKG_NAME)
  SECTION:=utils
  CATEGORY:=Utilities
  TITLE:=Gigaset Info Center (Go rewrite)
  DEPENDS:=+libc
endef

define Build/Compile
	cd $(PKG_BUILD_DIR) && go build -o $(PKG_BUILD_DIR)/$(PKG_NAME) .
endef

Package/$(PKG_NAME)/install:
	$(INSTALL_DIR) $(1:=$(StagingDir))/usr/bin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/$(PKG_NAME) $(1:=$(StagingDir))/usr/bin/

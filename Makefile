# Makefile for building gigaset-info-center Alpine/OpenWRT package
# Usage: make build / make install / make clean

PKG_NAME = gigaset-info-center
PKG_VER ?= 1.7

.PHONY: all build clean install help test

all: help

help:
	@echo "Usage:"
	@echo "  make build    - Build the .apk package"
	@echo "  make install  - Install the package on OpenWRT device"
	@echo "  make test     - Verify package structure"
	@echo "  make clean    - Clean build artifacts"
	@echo ""
	@echo "Environment variables:"
	@echo "  PKG_VER   - Package version (default: 1.7)"
	@echo "  HOST      - SSH host for OpenWRT device (default: root@192.168.1.1)"
	@echo "  DEST_DIR  - Destination directory on device (default: /root/packages)"

build:
	@bash scripts/build-package.sh "$(PKG_VER)"

install:
	@if [ ! -f dist/$(PKG_NAME)-$(PKG_VER)-r0.apk ]; then \
		echo "Error: Package not found. Run 'make build' first."; \
		exit 1; \
	fi
	@echo "Installing package on OpenWRT device..."
	@scp dist/$(PKG_NAME)-$(PKG_VER)-r0.apk $(HOST):$(DEST_DIR)/ || true
	@ssh $(HOST) "apk add --no-signature /root/packages/$(PKG_NAME)-$(PKG_VER)-r0.apk" || true
	@echo "Enabling service..."
	@ssh $(HOST) "rc-update add gigaset-info-center default" || true

test:
	@if [ ! -f dist/$(PKG_NAME)-$(PKG_VER)-r0.apk ]; then \
		echo "Error: Package not found. Run 'make build' first."; \
		exit 1; \
	fi
	@echo "Testing package structure..."
	@mkdir -p test-extract
	@tar xf dist/$(PKG_NAME)-$(PKG_VER)-r0.apk -C test-extract
	@test -f test-extract/gigaset-info-center/index.php || (echo "ERROR: index.php not found" && exit 1)
	@test -f test-extract/gigaset-info-center/weather.php || (echo "ERROR: weather.php not found" && exit 1)
	@test -f test-extract/gigaset-info-center/proxy.php || (echo "ERROR: proxy.php not found" && exit 1)
	@test -d test-extract/gigaset-info-center/icons || (echo "ERROR: icons directory not found" && exit 1)
	@echo "Package structure verified successfully!"

clean:
	@rm -rf build-stage dist test-extract

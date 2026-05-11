#!/bin/bash
# Build script for Alpine/OpenWRT .apk package
# Usage: ./scripts/build-package.sh [version]

set -e

PKG_NAME="gigaset-info-center"
VERSION="${1:-1.7}"
PKG_VER="${VERSION}-r0"

echo "=========================================="
echo "Building Alpine/OpenWRT Package (.apk)"
echo "Package: ${PKG_NAME}"
echo "Version: ${PKG_VER}"
echo "=========================================="

# Clean previous builds
rm -rf build-stage dist
mkdir -p dist

# Create staging directory structure
echo "Creating package structure..."
mkdir -p "build-stage/${PKG_NAME}/srv/gigaset-info-center/icons"
mkdir -p "build-stage/${PKG_NAME}/etc/lighttpd"
mkdir -p "build-stage/${PKG_NAME}/etc/init.d"
mkdir -p "build-stage/${PKG_NAME}/usr/share/doc/${PKG_NAME}"

# Copy web application files
echo "Copying web application files..."
cp index.php "build-stage/${PKG_NAME}/srv/gigaset-info-center/"
cp weather.php "build-stage/${PKG_NAME}/srv/gigaset-info-center/"
cp proxy.php "build-stage/${PKG_NAME}/srv/gigaset-info-center/"

# Copy icons
echo "Copying icons..."
for icon in icons/*.png; do
    cp "$icon" "build-stage/${PKG_NAME}/srv/gigaset-info-center/icons/"
done

# Copy configuration files
echo "Copying configuration files..."
cp etc/lighttpd/${PKG_NAME}.conf "build-stage/${PKG_NAME}/etc/lighttpd/"
cp etc/gigaset-env.example "build-stage/${PKG_NAME}/etc/"

# Copy init script and set permissions
echo "Setting up init script..."
cp ${PKG_NAME}.init "build-stage/${PKG_NAME}/etc/init.d/${PKG_NAME}"
chmod 755 "build-stage/${PKG_NAME}/etc/init.d/${PKG_NAME}"

# Copy documentation
echo "Copying documentation..."
cp LICENSE README.md "build-stage/${PKG_NAME}/usr/share/doc/${PKG_NAME}/"

# Create APKBUILD file (replacing PKG_NAME_VAR with actual name)
echo "Creating APKBUILD..."
sed 's/PKG_NAME_VAR/'"${PKG_NAME}"'/g' scripts/APKBUILD.template > build-stage/APKBUILD

# Create the .apk package using tar from parent directory
echo "Packaging..."
tar -cf dist/${PKG_NAME}-${PKG_VER}.apk -C build-stage ${PKG_NAME}

echo ""
echo "=========================================="
echo "Package created successfully!"
echo "File: dist/${PKG_NAME}-${PKG_VER}.apk"
echo "Size: $(du -sh dist/${PKG_NAME}-${PKG_VER}.apk | cut -f1)"
echo "=========================================="

# Show package contents for verification
echo ""
echo "Package contents:"
tar tf dist/${PKG_NAME}-${PKG_VER}.apk

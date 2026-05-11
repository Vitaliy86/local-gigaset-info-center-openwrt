# Replacement weather service for `info.gigaset.net`

This project allows to host your own mini-service for providing weather forecasts to your Gigaset IP handset (e.g. Gigaset C430A Go) even though, Siemens has discontinued running `info.gigaset.net` which was previously providing the data-source.

This replacement is currently extremely limited and only provides weather forecast data for a single location that you can configure. Configuration is done server-side and not via the handset (as previously).

This project is designed to run on ancient PHP5 because that is what my old trusty QNAP TS-109 Pro II is still on.

Of course, it would be really easy to reimplement this in `node`, `python` or whatever else, you would prefer.

Let me know what you think and send PRs for fixes and extensions!

## TL;DR

This project will:
- redirect `info.gigaset.net` to your own HTTP server on your local net
- provide weather data via two URLs that are used by the handset to retrieve "Info Center" data:
  - `http://info.gigaset.net/info/menu.jsp?lang=2&tz=120&mac=7C2F80XXXXXX&cc=49&handsetid=XXXXXXXXXX&provid=11`\
    This one is retrieved when accessing the "Info Center" from the "Extras" menu.
  - `http://info.gigaset.net/info/request.do?lang=2&tz=120&tick=true&mac=7C2F80XXXXXX&cc=49&provid=11&handsetid=XXXXXXXXXX`\
    This one is retrieved when data is retrieved for the standby screen scrolling text and when "Info Center" is used as the screen-saver. The `tick=true` parameter is sometimes present, sometimes not. I did not dig deeper.

## First step: Register a free account on OpenWeatherMap

Go to the OpenWeatherMap [signup](https://home.openweathermap.org/users/sign_up) page and create an account. Then, go to your [API keys](https://home.openweathermap.org/api_keys) page and copy your key for the next step.

## Prerequisite: Enable required PHP extensions

Before setting up the service, ensure the following PHP extensions are enabled on your server:

- **cURL** (`curl`) — Required to fetch weather data from the OpenWeatherMap API
- **GD** (`gd`) — Required to convert PNG weather icons to the Gigaset fnt format (bitmap)

On Synology DiskStation, you can enable these extensions in **Web Station** under the PHP settings for your virtual host.

## Second step: Set up the service

Copy the contents of this repository to a directory served as `http://<yourserver>/info`. In my case, this would be `/Qweb/info/` on the NAS.

To make your `apache` server run the provided scripts fine, add something like this to your `apache` configuration file, `/usr/local/apache/conf/apache.conf` in my case:

```apache
<Directory "/share/Qweb/info">
    DirectoryIndex menu.jsp
    Order deny,allow
    Deny from all
    Allow from localhost
    Allow from 192.168.10.0/24
	AddType application/x-httpd-php .jsp .do
	SetEnv OPENWEATHERMAP_API_KEY <key from step 1 here>
	SetEnv CITY "Berlin"
	SetEnv LATITUDE 52.52437
	SetEnv LONGITUDE 13.41053
</Directory>
```

Obviously, replace `192.168.10.0` with your local IP network and the `SetEnv` lines according to your situation.

Verify and activate your configuration:
```term
# /usr/local/apache/bin/apachectl configtest
# /usr/local/apache/bin/apachectl restart
```

Now, try and go to `http://<yourserver>/info`. You should see something like this:

<div style="overflow:scroll; height:12rem">
  <p style='text-align:center'>Do, 23.05.2024<br/>19,3/23,6°C/0 mm<br/>Bedeckt</p><p style='text-align:center'>Fr, 24.05.2024<br/>17,5/24,8°C/0 mm<br/>Bedeckt</p><p style='text-align:center'>Sa, 25.05.2024<br/>14,2/24,2°C/5 mm<br/>Leichter Regen/Bed.</p><p style='text-align:center'>So, 26.05.2024<br/>13,7/24,8°C/1 mm<br/>Leichter Regen/Mäßig bew.</p><p style='text-align:center'>Mo, 27.05.2024<br/>14,3/25,7°C/0 mm<br/>Bed./Überw. bew.</p><p style='text-align:center'>Di, 28.05.2024<br/>13,2/16,6°C/4 mm<br/>Leichter Regen</p>
</div>

Also, check `http://<yourserver>/info/menu.jsp` and `http://<yourserver>/info/request.do` which should show the same thing.

## Alternative: Synology DiskStation setup

If you are running a Synology DiskStation, you can use `.htaccess` instead of editing the global Apache config directly.

1. Copy the repository files to a web-served directory, e.g. `/volume1/web/info/`.
2. In **Web Station**, make sure PHP is enabled for the virtual host serving that directory. Also ensure the required PHP extensions are enabled:
   - **cURL** — for fetching OpenWeatherMap data
   - **GD** — for icon image processing
3. Copy `.htaccess.example` to `.htaccess` and fill in your values:

```apache
DirectoryIndex menu.php

Require local
Require ip 10.0.0.0/8

SetEnv OPENWEATHERMAP_API_KEY <your-api-key-here>
SetEnv CITY "Berlin"
SetEnv LATITUDE 52.52437
SetEnv LONGITUDE 13.41053

Options +FollowSymLinks
RewriteEngine On
RewriteRule ^menu\.jsp$ menu.php [L]
RewriteRule ^request\.do$ menu.php [L]
```

   Replace `10.0.0.0/8` with your local network range, and update the `SetEnv` lines with your city and coordinates.

   The `Require local` and `Require ip` directives restrict access to your local network only.

### Optional: Disable weather icons

By default, the service displays weather icons next to the forecast. If icons are not loading or you prefer a text-based view, add the following environment variable to your `.htaccess`:

```apache
SetEnv SHOW_ICONS false
```

When icons are disabled, a single German weather word is shown instead (e.g., "Regen", "Sonnig", "Nebel"). This provides a quick visual summary of the expected weather condition without requiring image support.

4. `.htaccess` is listed in `.gitignore` so your API key will not be accidentally committed to version control.

## Third step: Redirect your Gigaset phone to your own server

This very much depends on what kind of router you have. It is easy for routers that provide their own (caching) DNS service like e.g. OpenWRT. Also, if you can manually set up a mapping of host-names to IP addresses, you might get lucky. The key part is to make the DNS server that is configured in your Gigaset base station (probably via DHCP) to resolve:

    info.gigaset.net -> <your server IP>

In OpenWRT, you can set this up at `https://<your-router>/cgi-bin/luci/admin/network/dhcp` in the _Hostnames_ tab.


## Building for OpenWRT

### Build in OpenWRT SDK

Clone the OpenWRT SDK first:

```bash
# Download OpenWRT SDK for your target architecture
# Example for x86_64:
wget https://downloads.openwrt.org/releases/23.05.3/targets/x86/64/openwrt-sdk-23.05.3-x86-64_gcc-12.3.0_musl_x86_64.tar.xz
tar xf openwrt-sdk-23.05.3-x86_64_gcc-12.3.0_musl_x86_64.tar.xz
cd openwrt-sdk-23.05.3-x86_64_gcc-12.3.0_musl_x86_64

# Clone this repo into packages directory
git clone https://github.com/Vitaliy86/local-gigaset-info-center-openwrt.git package/gigaset-info-center
```

Then build the package:

```bash
# Make menuconfig and add gigaset-info-center under Network -> Other packages

# Build the package
make package/gigaset-info-center/compile V=s

# The .ipk package will be in:
# bin/packages/<architecture>/base/
```

### Install on OpenWRT 25.12.3+ device

OpenWRT 25.12.3+ uses Alpine Linux's `apk` package manager (opkg is deprecated).

```bash
# Download the .apk package to your OpenWRT device
wget -O /tmp/gigaset-info-center-1.7-r0.apk https://github.com/Vitaliy86/local-gigaset-info-center-openwrt/releases/download/v1.7/gigaset-info-center-1.7-r0.apk

# Install the package (skip signature check for locally built packages)
apk add --no-signature /tmp/gigaset-info-center-1.7-r0.apk

# Enable and start the service
rc-update add gigaset-info-center default
/etc/init.d/gigaset-info-center start
```

### Configure the service

After installation, edit the environment configuration file:

```bash
vi /etc/gigaset-info-center.env.example
# or copy to actual config:
cp /etc/gigaset-info-center.env.example /etc/gigaset-info-center.conf
```

Edit the following values:

| Variable | Description | Example |
|---|---|---|
| `OPENWEATHERMAP_API_KEY` | Your OpenWeatherMap API key | `a1b2c3d4e5f6...` |
| `CITY` | City name for display | `"Berlin"` |
| `LATITUDE` | Latitude coordinate | `52.52437` |
| `LONGITUDE` | Longitude coordinate | `13.41053` |
| `SHOW_ICONS` | Show weather icons (true/false) | `"true"` |

### Configure lighttpd

Add the gigaset-info-center config to your lighttpd configuration:

```bash
# Include the package config in /etc/lighttpd/lighttpd.conf
echo 'include "/etc/gigaset-info-center.conf"' >> /etc/lighttpd/lighttpd.conf

# Restart lighttpd
/etc/init.d/lighttpd restart
```

### Configure DNS redirect for Gigaset phone

In OpenWRT LuCI, set up DNS redirect:

1. Go to **Network > DHCP > DNS** tab
2. Add hostname alias: `info.gigaset.net` -> your router IP

Or via CLI:

```bash
# Add static DNS entry in /etc/config/dhcp
uci add dhcp static_host
uci set dhcp.static_host.ip="YOUR_ROUTER_IP"
uci set dhcp.static_host.name="info.gigaset.net"
uci commit dhcp
/etc/init.d.network restart
```

### Verify installation

Check if the service is running:

```bash
# Check service status
/etc/init.d/gigaset-info-center status

# Test locally (requires curl)
curl http://127.0.0.1:8081/
```

Expected output should show weather data in XHTML-GP format for your Gigaset phone.

## Build System

### File Structure

```
├── Makefile              # OpenWrt package build definition
├── gigaset-info-center.init # OpenWrt init script
├── etc/
│   ├── lighttpd/         # Lighttpd configuration
│   │   └── gigaset-info-center.conf
│   └── gigaset-env.example # Environment config template
├── icons/                # Weather icons
├── .github/workflows/    # GitHub Actions CI/CD
│   └── build-apk.yml     # Build workflow definition
```

### OpenWrt Package Makefile

The [`Makefile`](Makefile) follows OpenWrt package build system conventions:

```bash
# Build the package (run from OpenWrt SDK root)
make package/gigaset-info-center/compile V=s

# Clean build artifacts
make package/gigaset-info-center/clean

# Configure in menuconfig
make menuconfig
# Navigate to: Network -> Other packages -> gigaset-info-center
```

### Package Variables

| Variable | Description |
|---|---|
| `PKG_NAME` | Package name: `gigaset-info-center` |
| `PKG_VERSION` | Package version: `1.7` |
| `PKG_RELEASE` | Package release number |
| `DEPENDS` | Runtime dependencies |

### GitHub Actions Workflow

The [`.github/workflows/build-apk.yml`](.github/workflows/build-apk.yml) workflow:
- Builds `.apk` package on push to main/master or tag creation
- Uploads artifact for manual download
- Creates GitHub Release when tagging (`v*`)
- Tests package structure integrity

## References

I used information from

- https://www.ip-phone-forum.de/threads/gigaset-infodienst-selbst-gemacht.174719/ (thanks to [VoIPMaster](https://www.ip-phone-forum.de/members/voipmaster.95683/))
- https://copyandpastecode.blogspot.com/2008/08/siemens-s685ip-s68h.html (thanks to [Jon Bright](https://www.blogger.com/profile/13465823659620242219))
- and in particular:  http://www.ensued.net/request.do (coded in ruby)

See the latter for more ideas on how to include RSS updates and public transport information.

# gigaset-info-center (Go rewrite)

Replacement for the defunct `info.gigaset.net` weather service — delivers
3-day weather forecasts to Gigaset DECT IP phones running OpenWrt.

**v2.0**: Complete rewrite in Go. Single static binary, zero runtime
dependencies — replaces the old `php8 + lighttpd + php8-mod-gd +
php8-mod-curl + lighttpd-mod-fastcgi` stack (~20 MB installed) with one file
(~6 MB).

## How it works

```
Gigaset phone → DNS: info.gigaset.net → router IP
             → HTTP GET /info/menu.jsp
             → gigaset-info-center (Go binary on router)
             → OpenWeatherMap API
```

The binary also serves a PNG→FNT proxy at `/proxy/image.do` for weather icons.

## Installation

### 1. Copy APK to the router

```sh
scp gigaset-info-center-2.0-r0.apk root@192.168.1.1:/tmp/
```

### 2. Install (no internet required after download)

```sh
apk --allow-untrusted add /tmp/gigaset-info-center-2.0-r0.apk
```

### 3. Configure

```sh
cp /etc/gigaset-info-center.conf.example /etc/gigaset-info-center.conf
vi /etc/gigaset-info-center.conf
```

Minimum required settings:

```sh
OPENWEATHERMAP_API_KEY=abc123...   # free at openweathermap.org
CITY="Berlin"
LATITUDE=52.5200
LONGITUDE=13.4050
```

### 4. DNS: point info.gigaset.net to your router

Add to `/etc/config/dhcp` (OpenWrt dnsmasq):

```
config domain
    option name 'info.gigaset.net'
    option ip   '192.168.1.1'       # your router's LAN IP
```

Then reload: `service dnsmasq reload`

### 5. Port 80 (choose one option)

**Option A** — run directly on port 80 (default, simplest):

```sh
# If LuCI is on port 80, move it first:
uci set uhttpd.main.listen_http='0.0.0.0:8080'
uci commit uhttpd && service uhttpd restart
```

**Option B** — run on port 8080, redirect only phone traffic:

```sh
# In /etc/gigaset-info-center.conf:
LISTEN=:8080

# Redirect phone's requests on port 80 → 8080:
iptables -t nat -A PREROUTING \
    -p tcp --dport 80 \
    -s <phone-ip-or-subnet> \
    -j REDIRECT --to-port 8080
```

### 6. Start and enable autostart

```sh
service gigaset-info-center start
service gigaset-info-center enable
```

### 7. Test

```sh
curl -s http://localhost/info/menu.jsp
# Should return XHTML-GP with weather data
```

## Building from source

Requires Go 1.21+:

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -ldflags="-s -w" -o gigaset-info-center \
    ./cmd/gigaset-info-center/
```

Or use GitHub Actions — pushes to `main` build APKs for
`aarch64`, `armv7`, `x86_64`, `mips_sf` automatically.
Tag a commit (`git tag v2.0 && git push --tags`) to create a GitHub Release
with all APKs attached.

## Configuration reference

| Variable | Default | Description |
|---|---|---|
| `OPENWEATHERMAP_API_KEY` | *(required)* | OWM API key |
| `LATITUDE` | *(required)* | Decimal degrees |
| `LONGITUDE` | *(required)* | Decimal degrees |
| `CITY` | *(required)* | Shown in phone title bar |
| `LISTEN` | `:80` | TCP listen address |
| `SHOW_ICONS` | `true` | Show bitmap icons or text labels |
| `PROXY_BASE_URL` | `http://info.gigaset.net` | Server URL as seen by phone |
| `ICON_BASE_URL` | `https://openweathermap.org/img/wn` | PNG icon source |

## License

AGPL-3.0-or-later — see [LICENSE](LICENSE).
Original PHP version by Tilman Vogel.

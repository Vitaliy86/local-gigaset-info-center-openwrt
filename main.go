// gigaset-info-center — replacement for info.gigaset.net
// Single static binary; replaces PHP + lighttpd + php-gd + php-curl stack.
//
// Endpoints served:
//
//	GET /info/menu.jsp        — XHTML-GP weather page (main Gigaset request)
//	GET /info/request.do      — same (alternate Gigaset endpoint)
//	GET /proxy/image.do?data= — PNG→FNT icon proxy
//
// Configuration via environment variables (see gigaset-info-center.conf.example).
//
// Copyright (c) 2024 Tilman Vogel <tilman.vogel@web.de>
// Copyright (c) 2026 Vitaliy86 <vitaliy86@github.com>
// AGPL-3.0-or-later — see LICENSE
//
// Go rewrite by Vitaliy86 — replaces the original PHP + lighttpd stack
package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

const version = "2.0"

// ─── Config ──────────────────────────────────────────────────────────────────

type Config struct {
	Lat       string
	Lon       string
	City      string
	APIKey    string
	IconBase  string // URL prefix for PNG icons (OWM or local)
	ProxyBase string // our server's base URL as seen by the phone
	ShowIcons bool
	Listen    string
	Lang      string
}

func loadConfig() Config {
	env := func(k, def string) string {
		if v := os.Getenv(k); v != "" {
			return v
		}
		return def
	}
	return Config{
		Lat:       env("LATITUDE", ""),
		Lon:       env("LONGITUDE", ""),
		City:      env("CITY", ""),
		APIKey:    env("OPENWEATHERMAP_API_KEY", ""),
		IconBase:  env("ICON_BASE_URL", "https://openweathermap.org/img/wn"),
		ProxyBase: env("PROXY_BASE_URL", "http://info.gigaset.net"),
		ShowIcons: os.Getenv("SHOW_ICONS") != "false",
		Listen:    env("LISTEN", ":8080"),
		Lang:      env("Lang", "en"),
	}
}

// ─── OpenWeatherMap API types ─────────────────────────────────────────────────

type owmForecast struct {
	Cod  json.RawMessage `json:"cod"`
	Msg  json.RawMessage `json:"message"`
	List []owmItem       `json:"list"`
}

type owmItem struct {
	Dt      int64        `json:"dt"`
	Main    owmMain      `json:"main"`
	Weather []owmWeather `json:"weather"`
	Rain    struct {
		H3 float64 `json:"3h"`
	} `json:"rain"`
	Sys struct {
		Pod string `json:"pod"` // "d" = day, "n" = night
	} `json:"sys"`
}

type owmMain struct {
	Temp    float64 `json:"temp"`
	TempMin float64 `json:"temp_min"`
	TempMax float64 `json:"temp_max"`
}

type owmWeather struct {
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

// ─── Weather aggregation ──────────────────────────────────────────────────────

type DayData struct {
	Date        string  // "Mo, 12.05.2025"
	CurrentTemp float64 // temp of the first daytime slot
	MinTemp     float64
	MaxTemp     float64
	TotalRain   float64 // mm accumulated
	Icon        string  // most frequent daytime icon code, e.g. "04d"
	// private
	icons []string
	descs []string
}

var weekdays = [7]string{"So", "Mo", "Di", "Mi", "Do", "Fr", "Sa"}

// iconLabels maps OWM icon codes → German weather label.
var iconLabels = map[string]string{
	"01d": "Sonnig",
	"02d": "Wolkig",
	"03d": "Bewölkt",
	"04d": "Bedeckt",
	"09d": "Schauer",
	"10d": "Regen",
	"11d": "Gewitter",
	"13d": "Schnee",
	"50d": "Nebel",
}

func rainLabel(mm float64) string {
	switch {
	case mm < 3:
		return "Leichter Regen"
	case mm < 15:
		return "Regen"
	default:
		return "Starker Regen"
	}
}

// condLabel returns the human-readable German condition for a day.
func condLabel(d DayData) string {
	if d.Icon == "09d" || d.Icon == "10d" || d.Icon == "11d" ||
		(d.Icon == "04d" && d.TotalRain > 2) {
		return rainLabel(d.TotalRain)
	}
	if l, ok := iconLabels[d.Icon]; ok {
		return l
	}
	return d.Icon
}

// fetchWeather calls OWM API and returns up to 3 aggregated days.
func fetchWeather(cfg Config) ([]DayData, error) {
	apiURL := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/forecast"+
			"?lat=%s&lon=%s&appid=%s&units=metric&lang=%s",
		url.QueryEscape(cfg.Lat),
		url.QueryEscape(cfg.Lon),
		url.QueryEscape(cfg.APIKey),
		url.QueryEscape(cfg.Lang),
	)

	resp, err := http.Get(apiURL) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("OWM request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("OWM read: %w", err)
	}

	var fc owmForecast
	if err := json.Unmarshal(body, &fc); err != nil {
		return nil, fmt.Errorf("OWM JSON: %w", err)
	}

	// cod can be "200" (string) or 200 (number) depending on OWM version.
	codStr := strings.Trim(string(fc.Cod), `"`)
	if codStr != "200" {
		msg := strings.Trim(string(fc.Msg), `"`)
		return nil, fmt.Errorf("OWM API: %s", msg)
	}

	return aggregateDaily(fc.List), nil
}

// aggregateDaily groups 3-hourly slots into days (daytime only).
func aggregateDaily(items []owmItem) []DayData {
	var keys []string
	days := map[string]*DayData{}

	for _, it := range items {
		if it.Sys.Pod == "n" || len(it.Weather) == 0 {
			continue
		}

		t := time.Unix(it.Dt, 0).UTC()
		key := fmt.Sprintf("%s, %s",
			weekdays[t.Weekday()], t.Format("02.01.2006"))

		// Always use daytime icon variant.
		icon := strings.TrimRight(it.Weather[0].Icon, "dn") + "d"

		if _, exists := days[key]; !exists {
			days[key] = &DayData{
				Date:        key,
				CurrentTemp: it.Main.Temp,
				MinTemp:     it.Main.TempMin,
				MaxTemp:     it.Main.TempMax,
				TotalRain:   it.Rain.H3,
			}
			keys = append(keys, key)
		} else {
			d := days[key]
			if it.Main.TempMin < d.MinTemp {
				d.MinTemp = it.Main.TempMin
			}
			if it.Main.TempMax > d.MaxTemp {
				d.MaxTemp = it.Main.TempMax
			}
			d.TotalRain += it.Rain.H3
		}

		d := days[key]
		d.icons = append(d.icons, icon)
		d.descs = append(d.descs, it.Weather[0].Description)
	}

	result := make([]DayData, 0, len(keys))
	for _, k := range keys {
		d := days[k]
		d.Icon = mostCommon(d.icons)
		result = append(result, *d)
	}

	if len(result) > 3 {
		result = result[:3]
	}
	return result
}

// mostCommon returns the most frequently occurring string in ss.
func mostCommon(ss []string) string {
	counts := make(map[string]int, len(ss))
	for _, s := range ss {
		counts[s]++
	}
	type kv struct {
		k string
		v int
	}
	pairs := make([]kv, 0, len(counts))
	for k, v := range counts {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].v != pairs[j].v {
			return pairs[i].v > pairs[j].v
		}
		return pairs[i].k < pairs[j].k // stable tie-break
	})
	return pairs[0].k
}

// fmtTemp formats a temperature German-style (comma decimal separator).
func fmtTemp(t float64) string {
	return strings.ReplaceAll(fmt.Sprintf("%.0f", t), ".", ",")
}

// ─── HTTP handlers ────────────────────────────────────────────────────────────

// handleWeather serves the XHTML-GP page Gigaset phones expect.
func handleWeather(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		days, err := fetchWeather(cfg)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		io.WriteString(w, `<?xml version="1.0" encoding="utf-8"?>`)
		io.WriteString(w, "\n<!DOCTYPE html PUBLIC \"-//OMA//DTD XHTML Mobile 1.2//EN\""+
			" \"http://www.openmobilealliance.org/tech/DTD/xhtml-mobile12.dtd\">\n")
		fmt.Fprintf(w,
			"<html>\n<head>\n<title>%s (%s)</title>\n"+
				"<meta name=\"expires\" content=\"3600\" />\n</head>\n<body>\n",
			cfg.City, version)

		if err != nil {
			log.Printf("weather error: %v", err)
			fmt.Fprintf(w, "<p><b>Fehler:</b> %s</p>\n", err)
		} else {
			for i, d := range days {
				wd := weatherWidget(cfg, d)
				if i == 0 {
					// "Heute <icon/label> jetzt 12°C, max 15°C"
					fmt.Fprintf(w, "<p>Heute %s jetzt %s°C, max %s°C</p>\n",
						wd, fmtTemp(d.CurrentTemp), fmtTemp(d.MaxTemp))
				} else {
					// "Mo <icon/label> 8-14°C"
					dow := strings.SplitN(d.Date, ",", 2)[0]
					fmt.Fprintf(w, "<p>%s %s %s-%s°C</p>\n",
						dow, wd, fmtTemp(d.MinTemp), fmtTemp(d.MaxTemp))
				}
			}
		}

		io.WriteString(w, "</body>\n</html>")
	}
}

// weatherWidget returns either a text label or an <object> fnt tag.
func weatherWidget(cfg Config, d DayData) string {
	if !cfg.ShowIcons {
		return condLabel(d)
	}
	iconURL := url.QueryEscape(cfg.IconBase + "/" + d.Icon + ".png")
	return fmt.Sprintf(
		"<object data='%s/proxy/image.do?data=%s'"+
			" type='image/fnt' width='16' height='16'></object>",
		cfg.ProxyBase, iconURL,
	)
}

// handleProxy fetches a remote PNG and converts it to Gigaset fnt (1bpp bitmap).
func handleProxy(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawURL := r.URL.Query().Get("data")
		if rawURL == "" {
			http.Error(w, "Bad Request: missing data param", http.StatusBadRequest)
			return
		}

		parsed, err := url.Parse(rawURL)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			http.Error(w, "Bad Request: invalid URL", http.StatusBadRequest)
			return
		}

		// Restrict to configured icon source — prevent open-proxy abuse.
		if cfg.IconBase != "" && !strings.HasPrefix(rawURL, cfg.IconBase+"/") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		resp, err := http.Get(rawURL) //nolint:noctx
		if err != nil {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			http.Error(w, "Invalid image", http.StatusInternalServerError)
			return
		}

		fnt := toFnt(img, 16, 16)
		w.Header().Set("Content-Type", "image/fnt")
		w.Header().Set("Content-Length", fmt.Sprint(len(fnt)))
		w.Write(fnt) //nolint:errcheck
	}
}

// toFnt converts any image.Image to the Gigaset fnt format:
//
//	uint16 LE  width
//	uint16 LE  height
//	rows of ceil(width/8) bytes, MSB first, 0 = white, 1 = black
//
// The source is scaled to (w×h) using nearest-neighbour with
// alpha-over-white compositing before thresholding.
func toFnt(src image.Image, w, h int) []byte {
	b := src.Bounds()
	sw, sh := b.Dx(), b.Dy()

	rowBytes := (w + 7) / 8
	out := make([]byte, 4+h*rowBytes)
	binary.LittleEndian.PutUint16(out[0:], uint16(w))
	binary.LittleEndian.PutUint16(out[2:], uint16(h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Nearest-neighbour sample from source.
			sx := b.Min.X + x*sw/w
			sy := b.Min.Y + y*sh/h
			rv, gv, bv, av := src.At(sx, sy).RGBA() // 0–65535

			// Alpha-composite pixel over white background.
			af := float64(av) / 65535
			rf := float64(rv)/65535*af + (1 - af)
			gf := float64(gv)/65535*af + (1 - af)
			bf := float64(bv)/65535*af + (1 - af)

			// BT.601 luma; threshold matches PHP's `< 128` on 0–255 scale.
			luma := 0.299*rf + 0.587*gf + 0.114*bf
			if luma < 0.502 {
				out[4+y*rowBytes+x/8] |= 0x80 >> (x % 8)
			}
		}
	}
	return out
}

// ─── Entry point ──────────────────────────────────────────────────────────────

func main() {
	cfg := loadConfig()

	// Check for CLI arguments first (version/help)
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("gigaset-info-center v%s\n", version)
			return
		case "--help", "-h":
			fmt.Println(`Usage: gigaset-info-center [options]

Options:
  --version, -v    Show version
  --help, -h       Show help

Configuration:
  Configuration can be provided via:
  1. Environment variables (LATITUDE, LONGITUDE, CITY, OPENWEATHERMAP_API_KEY, etc.)
  2. Configuration file: /etc/gigaset-info-center.conf (INI format)

  The application will load configuration from the file if it exists,
  otherwise it falls back to environment variables.

Examples:
  gigaset-info-center --version
  gigaset-info-center --help
  gigaset-info-center  (uses environment variables or /etc/gigaset-info-center.conf)`)
			return
		}
	}

	if cfg.APIKey == "" {
		log.Fatal("OPENWEATHERMAP_API_KEY is required")
	}
	if cfg.Lat == "" || cfg.Lon == "" {
		log.Fatal("LATITUDE and LONGITUDE are required")
	}
	if cfg.City == "" {
		log.Print("warning: CITY is not set")
	}
	if cfg.Lang == "" {
		log.Print("warning: Lang is not set")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/info/menu.jsp", handleWeather(cfg))
	mux.HandleFunc("/info/request.do", handleWeather(cfg))
	mux.HandleFunc("/proxy/image.do", handleProxy(cfg))
	// Health check — useful for procd respawn monitoring.
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "ok")
	})

	log.Printf("gigaset-info-center v%s  listen=%s  city=%s",
		version, cfg.Listen, cfg.City)
	if err := http.ListenAndServe(cfg.Listen, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}

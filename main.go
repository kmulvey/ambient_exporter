package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"go.szostok.io/version"
	"go.szostok.io/version/printer"
)

type WeatherReport struct {
	// Station info
	PassKey     string    `json:"passkey"`
	StationType string    `json:"stationtype"`
	DateUTC     time.Time `json:"dateutc"`

	// Outdoor conditions
	Temp           float64 `json:"tempf"`
	Humidity       int     `json:"humidity"`
	WindSpeedMph   float64 `json:"windspeedmph"`
	WindGustMph    float64 `json:"windgustmph"`
	MaxDailyGust   float64 `json:"maxdailygust"`
	WindDir        int     `json:"winddir"`
	UV             int     `json:"uv"`
	SolarRadiation float64 `json:"solarradiation"`
	BattOut        int     `json:"battout"`

	// Calculated
	AtmosphericPressure  float64 `json:"atmosphericpressure"`
	DewPoint             float64 `json:"dewpoint"`
	HeatIndex            float64 `json:"heatindex"`
	WindChill            float64 `json:"windchill"`
	AbsoluteHumidity     float64 `json:"absolutehumidity"`
	VaporPressure        float64 `json:"vaporpressure"`
	VaporPressureDeficit float64 `json:"vaporpressuredeficit"`
	WindRun              float64 `json:"windrun"`
	Evapotranspiration   float64 `json:"evapotranspiration"`

	// Rain
	HourlyRainIn  float64 `json:"hourlyrainin"`
	EventRainIn   float64 `json:"eventrainin"`
	DailyRainIn   float64 `json:"dailyrainin"`
	WeeklyRainIn  float64 `json:"weeklyrainin"`
	MonthlyRainIn float64 `json:"monthlyrainin"`
	YearlyRainIn  float64 `json:"yearlyrainin"`
	TotalRainIn   float64 `json:"totalrainin"`

	// Indoor conditions
	TempIn     float64 `json:"tempinf"`
	HumidityIn int     `json:"humidityin"`

	// Barometric pressure
	BaromRelIn float64 `json:"baromrelin"`
	BaromAbsIn float64 `json:"baromabsin"`
}

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 15:04:05",
	})

	var addr string
	var elevationFt float64
	var v bool

	flag.StringVar(&addr, "addr", ":9600", "address for the http server")
	flag.Float64Var(&elevationFt, "elevation", 0, "station elevation in feet")
	flag.BoolVar(&v, "version", false, "print version")
	flag.BoolVar(&v, "v", false, "print version")
	flag.Parse()

	if elevationFt == 0 {
		log.Fatal("elevation is required")
	}

	if v {
		var verPrinter = printer.New()
		var info = version.Get()
		if err := verPrinter.PrintInfo(os.Stdout, info); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		log.Infof("HTTP server starting on %s", addr)
		srv := &http.Server{
			Addr:         addr,
			Handler:      nil,
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
		}
		http.Handle("/metrics", promhttp.Handler())
		http.HandleFunc("/data/report/", func(w http.ResponseWriter, r *http.Request) {
			reportHandler(w, r, elevationFt)
		})
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-sigChan
	os.Exit(0)
}

func reportHandler(w http.ResponseWriter, r *http.Request, elevationFt float64) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle malformed query string that starts with & instead of ?
	rawQuery := r.URL.RawQuery
	if rawQuery == "" && strings.Contains(r.URL.Path, "&") {
		// Extract query from path
		parts := strings.SplitN(r.URL.Path, "&", 2)
		if len(parts) == 2 {
			rawQuery = parts[1]
		}
	}

	query, err := url.ParseQuery(rawQuery)
	if err != nil {
		log.Errorf("Failed to parse query: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	report, err := parseWeatherReport(query)
	if err != nil {
		log.Errorf("Failed to parse weather report: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Calculate derived values
	calculateDerivedMeasures(report, elevationFt)

	updateMetrics(report)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "received"}); err != nil {
		log.Errorf("Failed to encode response: %v", err)
	}
}

func parseWeatherReport(query url.Values) (*WeatherReport, error) {
	report := &WeatherReport{
		PassKey:     query.Get("PASSKEY"),
		StationType: query.Get("stationtype"),
	}

	// Parse timestamp
	if dateStr := query.Get("dateutc"); dateStr != "" {
		// Replace + with space for parsing
		dateStr = strings.ReplaceAll(dateStr, "+", " ")
		t, err := time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		report.DateUTC = t
	}

	// Outdoor conditions
	report.Temp = parseTempF("tempf", query)
	report.Humidity = parseInt("humidity", query)
	report.WindSpeedMph = parseFloat("windspeedmph", query)
	report.WindGustMph = parseFloat("windgustmph", query)
	report.MaxDailyGust = parseFloat("maxdailygust", query)
	report.WindDir = parseInt("winddir", query)
	report.UV = parseInt("uv", query)
	report.SolarRadiation = parseFloat("solarradiation", query)
	report.BattOut = parseInt("battout", query)

	// Rain
	report.HourlyRainIn = parseFloat("hourlyrainin", query)
	report.EventRainIn = parseFloat("eventrainin", query)
	report.DailyRainIn = parseFloat("dailyrainin", query)
	report.WeeklyRainIn = parseFloat("weeklyrainin", query)
	report.MonthlyRainIn = parseFloat("monthlyrainin", query)
	report.YearlyRainIn = parseFloat("yearlyrainin", query)
	report.TotalRainIn = parseFloat("totalrainin", query)

	// Indoor conditions
	report.TempIn = parseTempF("tempinf", query)
	report.HumidityIn = parseInt("humidityin", query)

	// Barometric pressure
	report.BaromRelIn = parseFloat("baromrelin", query)
	report.BaromAbsIn = parseFloat("baromabsin", query)

	return report, nil
}

func parseFloat(key string, query url.Values) float64 {
	if val := query.Get(key); val != "" {
		f, _ := strconv.ParseFloat(val, 64)
		return f
	}
	return 0
}

func parseTempF(key string, query url.Values) float64 {
	f := parseFloat(key, query)
	if f == 0 {
		return 0
	}
	return (f - 32) * 5 / 9
}

func parseInt(key string, query url.Values) int {
	if val := query.Get(key); val != "" {
		i, _ := strconv.Atoi(val)
		return i
	}
	return 0
}

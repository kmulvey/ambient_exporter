# ambient_exporter

`ambient_exporter` is a Prometheus exporter for [Ambient Weather stations](https://ambientweather.com/). It receives weather data directly from your station (via the "Customized" server feature) and exposes it as Prometheus metrics.

In addition to the raw data provided by the station, this exporter calculates several derived metrics such as Dew Point, Heat Index, Wind Chill, and Evapotranspiration.

## Features

-   **Direct Ingestion**: Listens for HTTP GET requests from your weather station.
-   **Prometheus Metrics**: Exposes weather data at `/metrics`.
-   **Derived Metrics**: Automatically calculates:
    -   Atmospheric Pressure (adjusted for elevation)
    -   Dew Point
    -   Heat Index
    -   Wind Chill
    -   Absolute Humidity
    -   Vapor Pressure & Vapor Pressure Deficit
    -   Wind Run
    -   Evapotranspiration
-   **Unit Conversion**: Converts temperature to Celsius for standard metric export (raw inputs are typically Fahrenheit).

## Installation

### From Source

```bash
go install github.com/kmulvey/ambient_exporter@latest
```

### Build Locally

```bash
git clone https://github.com/kmulvey/ambient_exporter.git
cd ambient_exporter
go build
```

## Usage

The exporter requires your station's elevation (in feet) to accurately calculate atmospheric pressure and other derived metrics.

```bash
./ambient_exporter -elevation <ELEVATION_IN_FEET>
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-addr` | Address for the HTTP server to listen on | `:9600` |
| `-elevation` | Station elevation in feet (**Required**) | `0` |
| `-wind-dir-offset-deg` | Degrees clockwise to add to the reported wind direction (use to correct station mounting alignment; may be negative) | `0` |
| `-version` | Print version information | `false` |

## Weather Station Configuration

To send data to this exporter, you need to configure your Ambient Weather station (e.g., WS-2902, WS-2000, WS-5000) to upload to a "Customized" server.

1.  Open your weather station's configuration app (e.g., awnet) or web interface.
2.  Navigate to the **Upload** or **Customized** server settings.
3.  Enter the following details:
    -   **Server IP / Hostname**: The IP address or hostname of the machine running `ambient_exporter`.
    -   **Path**: `/data/report/`
    -   **Port**: `9600` (or whatever you configured with `-addr`)
    -   **Protocol**: HTTP
    -   **Upload Interval**: As desired (e.g., 16 or 60 seconds).

## Metrics

Metrics are exposed at `http://<host>:9600/metrics`.

### Exported Metrics

-   `ambient_weather_station_info`: Station type and passkey.
-   `ambient_weather_outdoor_temperature_celsius`
-   `ambient_weather_outdoor_humidity_percent`
-   `ambient_weather_wind_speed_mph`
-   `ambient_weather_wind_gust_mph`
-   `ambient_weather_max_daily_gust_mph`
-   `ambient_weather_wind_dir_degrees`
-   `ambient_weather_uv_index`
-   `ambient_weather_solar_radiation_wm2`
-   `ambient_weather_battery_ok`: 1 if battery is OK, 0 otherwise.
-   `ambient_weather_rain_hourly_in`
-   `ambient_weather_rain_event_in`
-   `ambient_weather_rain_daily_in`
-   `ambient_weather_rain_weekly_in`
-   `ambient_weather_rain_monthly_in`
-   `ambient_weather_rain_yearly_in`
-   `ambient_weather_rain_total_in`
-   `ambient_weather_indoor_temperature_celsius`
-   `ambient_weather_indoor_humidity_percent`
-   `ambient_weather_barom_rel_in`
-   `ambient_weather_barom_abs_in`

### Derived Metrics

-   `ambient_weather_atmospheric_pressure_hpa`
-   `ambient_weather_dew_point_celsius`
-   `ambient_weather_heat_index_celsius`
-   `ambient_weather_wind_chill_celsius`
-   `ambient_weather_absolute_humidity_gm3`
-   `ambient_weather_vapor_pressure_kpa`
-   `ambient_weather_vapor_pressure_deficit_kpa`
-   `ambient_weather_wind_run_miles`
-   `ambient_weather_evapotranspiration_in`

## Systemd Service

A sample systemd service file is provided in `ambient_exporter.service`.

1.  Create the configuration file `/etc/default/ambient_exporter`. This file is required to set the elevation:
    ```bash
    AMBIENT_ELEVATION=1234
    # Optional: Override default address
    # AMBIENT_ADDR=:9600
    ```
2.  Edit `ambient_exporter.service` to ensure `ExecStart` points to the correct location of the binary (default is `/usr/bin/ambient_exporter`).
3.  Copy to systemd directory:
    ```bash
    sudo cp ambient_exporter.service /etc/systemd/system/
    ```
4.  Reload and start:
    ```bash
    sudo systemctl daemon-reload
    sudo systemctl enable --now ambient_exporter
    ```

## Grafana Dashboard

A Grafana dashboard configuration is included in `grafana-config.json`. You can import this JSON file directly into Grafana to visualize your weather data.

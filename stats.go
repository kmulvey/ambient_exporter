package main

import "github.com/prometheus/client_golang/prometheus"

var (
	// Station info
	stationInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ambient_weather_station_info",
			Help: "Ambient Weather station information",
		},
		[]string{"station_type"},
	)

	// Outdoor temperature
	outdoorTemp = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_outdoor_temperature_celsius",
			Help: "Outdoor temperature in Celsius",
		},
	)

	// Outdoor humidity
	outdoorHumidity = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_outdoor_humidity_percent",
			Help: "Outdoor relative humidity percentage",
		},
	)

	// Wind speed
	windSpeedMph = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_wind_speed_mph",
			Help: "Average wind speed in mph",
		},
	)

	// Wind gust
	windGustMph = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_wind_gust_mph",
			Help: "Current wind gust in mph",
		},
	)

	// Max daily gust
	maxDailyGust = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_max_daily_gust_mph",
			Help: "Maximum wind gust today in mph",
		},
	)

	// Wind direction
	windDir = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_wind_direction_degrees",
			Help: "Wind direction in degrees (0-360)",
		},
	)

	// UV index
	uvIndex = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_uv_index",
			Help: "UV index",
		},
	)

	// Solar radiation
	solarRadiation = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_solar_radiation_wm2",
			Help: "Solar radiation in W/m²",
		},
	)

	// Battery status
	batteryOut = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_outdoor_battery_ok",
			Help: "Outdoor sensor battery status (1 = OK)",
		},
	)

	// Rain - hourly
	hourlyRainIn = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_rain_hourly_inches",
			Help: "Rain in last hour in inches",
		},
	)

	// Rain - event
	eventRainIn = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_rain_event_inches",
			Help: "Rain since last reset in inches",
		},
	)

	// Rain - daily
	dailyRainIn = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_rain_daily_inches",
			Help: "Rain today in inches",
		},
	)

	// Rain - weekly
	weeklyRainIn = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_rain_weekly_inches",
			Help: "Rain this week in inches",
		},
	)

	// Rain - monthly
	monthlyRainIn = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_rain_monthly_inches",
			Help: "Rain this month in inches",
		},
	)

	// Rain - yearly
	yearlyRainIn = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_rain_yearly_inches",
			Help: "Rain this year in inches",
		},
	)

	// Rain - total
	totalRainIn = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_rain_total_inches",
			Help: "Lifetime rain total in inches",
		},
	)

	// Indoor temperature
	indoorTemp = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_indoor_temperature_celsius",
			Help: "Indoor temperature in Celsius",
		},
	)

	// Indoor humidity
	indoorHumidity = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_indoor_humidity_percent",
			Help: "Indoor relative humidity percentage",
		},
	)

	// Barometric pressure - relative
	baromRelIn = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_barometric_pressure_relative_inhg",
			Help: "Sea-level (relative) barometric pressure in inHg",
		},
	)

	// Barometric pressure - absolute
	baromAbsIn = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_barometric_pressure_absolute_inhg",
			Help: "Absolute station barometric pressure in inHg",
		},
	)

	// Calculated fields
	atmosphericPressure = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_atmospheric_pressure_hpa",
			Help: "Atmospheric pressure at station elevation in hPa",
		},
	)

	dewPoint = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_dew_point_celsius",
			Help: "Dew point in Celsius",
		},
	)

	heatIndex = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_heat_index_celsius",
			Help: "Heat index in Celsius",
		},
	)

	windChill = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_wind_chill_celsius",
			Help: "Wind chill in Celsius",
		},
	)

	absoluteHumidity = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_absolute_humidity_gm3",
			Help: "Absolute humidity in g/m³",
		},
	)

	vaporPressure = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_vapor_pressure_hpa",
			Help: "Vapor pressure in hPa",
		},
	)

	vaporPressureDeficit = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_vapor_pressure_deficit_hpa",
			Help: "Vapor pressure deficit in hPa",
		},
	)

	windRun = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_wind_run_miles",
			Help: "Wind run in miles (24-hour period)",
		},
	)

	evapotranspiration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ambient_weather_evapotranspiration_mm_day",
			Help: "Evapotranspiration in mm/day",
		},
	)
)

func init() {
	prometheus.MustRegister(
		stationInfo,
		outdoorTemp,
		outdoorHumidity,
		windSpeedMph,
		windGustMph,
		maxDailyGust,
		windDir,
		uvIndex,
		solarRadiation,
		batteryOut,
		hourlyRainIn,
		eventRainIn,
		dailyRainIn,
		weeklyRainIn,
		monthlyRainIn,
		yearlyRainIn,
		totalRainIn,
		indoorTemp,
		indoorHumidity,
		baromRelIn,
		baromAbsIn,
		atmosphericPressure,
		dewPoint,
		heatIndex,
		windChill,
		absoluteHumidity,
		vaporPressure,
		vaporPressureDeficit,
		windRun,
		evapotranspiration,
	)
}

// updateMetrics updates all Prometheus metrics from the weather report
func updateMetrics(r *WeatherReport) {
	stationInfo.WithLabelValues(r.StationType).Set(1)
	outdoorTemp.Set(r.Temp)
	outdoorHumidity.Set(float64(r.Humidity))
	windSpeedMph.Set(r.WindSpeedMph)
	windGustMph.Set(r.WindGustMph)
	maxDailyGust.Set(r.MaxDailyGust)
	windDir.Set(r.WindDir)
	uvIndex.Set(float64(r.UV))
	solarRadiation.Set(r.SolarRadiation)
	batteryOut.Set(float64(r.BattOut))
	hourlyRainIn.Set(r.HourlyRainIn)
	eventRainIn.Set(r.EventRainIn)
	dailyRainIn.Set(r.DailyRainIn)
	weeklyRainIn.Set(r.WeeklyRainIn)
	monthlyRainIn.Set(r.MonthlyRainIn)
	yearlyRainIn.Set(r.YearlyRainIn)
	totalRainIn.Set(r.TotalRainIn)
	indoorTemp.Set(r.TempIn)
	indoorHumidity.Set(float64(r.HumidityIn))
	baromRelIn.Set(r.BaromRelIn)
	baromAbsIn.Set(r.BaromAbsIn)
	atmosphericPressure.Set(r.AtmosphericPressure)
	dewPoint.Set(r.DewPoint)
	heatIndex.Set(r.HeatIndex)
	windChill.Set(r.WindChill)
	absoluteHumidity.Set(r.AbsoluteHumidity)
	vaporPressure.Set(r.VaporPressure)
	vaporPressureDeficit.Set(r.VaporPressureDeficit)
	windRun.Set(r.WindRun)
	evapotranspiration.Set(r.Evapotranspiration)
}

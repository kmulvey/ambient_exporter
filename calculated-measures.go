package main

import "math"

func calculateDerivedMeasures(report *WeatherReport, elevationFt float64) {
	report.AtmosphericPressure = calculateAtmosphericPressure(elevationFt)
	report.DewPoint = calculateDewPoint(report.Temp, report.Humidity)
	report.HeatIndex = calculateHeatIndex(report.Temp, report.Humidity)
	report.WindChill = calculateWindChill(report.Temp, report.WindSpeedMph)
	report.AbsoluteHumidity = calculateAbsoluteHumidity(report.Temp, report.Humidity, elevationFt)
	report.VaporPressure = calculateVaporPressure(report.Temp, report.Humidity)
	report.VaporPressureDeficit = calculateVaporPressureDeficit(report.Temp, report.Humidity)
	report.WindRun = calculateWindRun(report.WindSpeedMph)
	report.Evapotranspiration = calculateEvapotranspiration(report.Temp, report.Humidity, report.WindSpeedMph, report.SolarRadiation, elevationFt)
}

// calculateAtmosphericPressure calculates pressure at elevation in hPa
func calculateAtmosphericPressure(elevationFt float64) float64 {
	elevationM := elevationFt * 0.3048
	// Standard atmospheric pressure at sea level
	p0 := 1013.25 // hPa
	// Barometric formula
	pressure := p0 * math.Pow(1-(0.0065*elevationM/288.15), 5.255)
	return pressure
}

func calculateDewPoint(tempC float64, humidity int) float64 {
	// Magnus formula (not significantly affected by elevation)
	a := 17.27
	b := 237.7
	alpha := ((a * tempC) / (b + tempC)) + math.Log(float64(humidity)/100.0)
	dewPoint := (b * alpha) / (a - alpha)
	return dewPoint
}

func calculateHeatIndex(tempC float64, humidity int) float64 {
	// Convert Celsius to Fahrenheit
	tempF := (tempC * 9.0 / 5.0) + 32.0

	// Heat index formula (less accurate at high altitude, but commonly used)
	// Note: Some sources suggest heat index shouldn't be used above ~3000ft
	hiF := -42.379 + 2.04901523*tempF + 10.14333127*float64(humidity) - 0.22475541*tempF*float64(humidity) -
		6.83783*math.Pow(10, -3)*math.Pow(tempF, 2) - 5.481717*math.Pow(10, -2)*math.Pow(float64(humidity), 2) +
		1.22874*math.Pow(10, -3)*math.Pow(tempF, 2)*float64(humidity) + 8.5282*math.Pow(10, -4)*tempF*math.Pow(float64(humidity), 2) -
		1.99*math.Pow(10, -6)*math.Pow(tempF, 2)*math.Pow(float64(humidity), 2)

	// Convert back to Celsius
	hiC := (hiF - 32.0) * 5.0 / 9.0
	return hiC
}

func calculateWindChill(tempC float64, windSpeedMph float64) float64 {
	// Convert Celsius to Fahrenheit
	tempF := (tempC * 9.0 / 5.0) + 32.0

	// Wind chill formula (altitude doesn't significantly affect this)
	wcF := 35.74 + 0.6215*tempF - 35.75*math.Pow(windSpeedMph, 0.16) + 0.4275*tempF*math.Pow(windSpeedMph, 0.16)

	// Convert back to Celsius
	wcC := (wcF - 32.0) * 5.0 / 9.0
	return wcC
}

func calculateAbsoluteHumidity(tempC float64, humidity int, elevationFt float64) float64 {
	// Calculate saturation vapor pressure (in hPa)
	svp := 6.112 * math.Exp((17.67*tempC)/(tempC+243.5))
	// Calculate actual vapor pressure
	avp := svp * float64(humidity) / 100.0

	// Adjust for altitude - air density is lower
	atmosphericPressure := calculateAtmosphericPressure(elevationFt)
	densityRatio := atmosphericPressure / 1013.25

	// Calculate absolute humidity (in g/m³) with altitude correction
	ah := (216.7 * avp) / (tempC + 273.15) * densityRatio
	return ah
}

func calculateVaporPressure(tempC float64, humidity int) float64 {
	// Calculate saturation vapor pressure (in hPa)
	svp := 6.112 * math.Exp((17.67*tempC)/(tempC+243.5))
	// Calculate actual vapor pressure
	avp := svp * float64(humidity) / 100.0
	return avp
}

func calculateVaporPressureDeficit(tempC float64, humidity int) float64 {
	// Calculate saturation vapor pressure (in hPa)
	svp := 6.112 * math.Exp((17.67*tempC)/(tempC+243.5))
	// Calculate actual vapor pressure
	avp := svp * float64(humidity) / 100.0
	// Calculate vapor pressure deficit
	vpd := svp - avp
	return vpd
}

func calculateWindRun(windSpeedMph float64) float64 {
	// Wind run is typically calculated over a period of time.
	// Here we assume a 24-hour period for demonstration purposes.
	hours := 24.0
	windRun := windSpeedMph * hours
	return windRun
}

func calculateEvapotranspiration(tempC float64, humidity int, windSpeedMph float64, solarRadiation float64, elevationFt float64) float64 {
	// Simplified FAO Penman-Monteith equation adjusted for elevation
	elevationM := elevationFt * 0.3048

	// Convert wind speed from mph to m/s
	windSpeedMs := windSpeedMph * 0.44704

	// Adjust wind speed for elevation (wind is typically stronger at higher elevations)
	// Reference height adjustment factor
	windSpeedMs *= math.Pow(elevationM/10.0, 0.2)

	// Get atmospheric pressure at elevation
	atmosphericPressure := calculateAtmosphericPressure(elevationFt)

	// Psychrometric constant adjusted for elevation (kPa/°C)
	gamma := 0.665 * math.Pow(10, -3) * atmosphericPressure

	// Slope of saturation vapor pressure curve
	delta := 4098 * (0.6108 * math.Exp((17.27*tempC)/(tempC+237.3))) / math.Pow(tempC+237.3, 2)

	// Saturation vapor pressure
	es := 0.6108 * math.Exp((17.27*tempC)/(tempC+237.3))

	// Actual vapor pressure
	ea := es * float64(humidity) / 100.0

	// Calculate ET (mm/day) with altitude corrections
	numerator := 0.408*delta*solarRadiation + gamma*(900/(tempC+273))*windSpeedMs*(es-ea)
	denominator := delta + gamma*(1+0.34*windSpeedMs)
	et := numerator / denominator

	return et
}

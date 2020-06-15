// Copyright (c) 2020 Tanner Ryan. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

package parser

import "encoding/json"

// ConditionsHTTP is for parsing weather conditions retrieved over HTTP.
type ConditionsHTTP struct {
	Data  *WeatherHTTP `json:"data"`  // Data is weather conditions from HTTP
	Error *Error       `json:"error"` // Error is a general error message
}

// WeatherHTTP is for parsing HTTP weather data.
type WeatherHTTP struct {
	DeviceID   string      `json:"did"`        // DeviceID is unique device ID
	Timestamp  *Time       `json:"ts"`         // Timestamp is device timestamp (epoch seconds)
	Conditions []EntryHTTP `json:"conditions"` // Conditions is a list of weather conditions
}

// EntryHTTP is a union type for containing multiple types of weather
// conditions.
type EntryHTTP struct {
	LogicalSensorID   *int        `json:"lsid"`                // LogicalSensorID is logical sensor ID
	DataStructureType *RecordType `json:"data_structure_type"` // DataStructureType indicates data structure
	Values            interface{} // Values are the weather condition values (generated)
}

// conditionHeader is used for partial weather conditions parsing.
type conditionHeader struct {
	LogicalSensorID   *int        `json:"lsid"`                // LogicalSensorID is logical sensor ID
	DataStructureType *RecordType `json:"data_structure_type"` // DataStructureType indicates data structure
}

// WeatherISS is weather conditions from ISS.
type WeatherISS struct {
	TransmitterID *int `json:"txid"` // TransmitterID is ID of transmitter

	Temperature *float64 `json:"temp"`       // Temperature (°F)
	Humidity    *float64 `json:"hum"`        // Humidity (%RH)
	Dewpoint    *float64 `json:"dew_point"`  // Dewpoint (°F)
	Wetbulb     *float64 `json:"wet_bulb"`   // Wetbulb (°F)
	HeatIndex   *float64 `json:"heat_index"` // HeatIndex (°F)
	WindChill   *float64 `json:"wind_chill"` // WindChill (°F)
	THWIndex    *float64 `json:"thw_index"`  // THWIndex is "feels like" (°F)
	THSWIndex   *float64 `json:"thsw_index"` // THWSIndex is "feels like" including solar (°F)

	WindSpeedLast          *float64 `json:"wind_speed_last"`                  // WindSpeedLast is most recent wind speed (mph)
	WindDirLast            *float64 `json:"wind_dir_last"`                    // WindDirLast is most recent wind direction (°)
	WindSpeedAvgLast1Min   *float64 `json:"wind_speed_avg_last_1_min"`        // WindSpeedAvgLast1Min is average wind over last minute (mph)
	WindDirAvgLast1Min     *float64 `json:"wind_dir_scalar_avg_last_1_min"`   // WindDirAvgLast1Min is average wind direction over last minute (°)
	WindSpeedAvgLast2Min   *float64 `json:"wind_speed_avg_last_2_min"`        // WindSpeedAvgLast2Min is average wind over last 2 minutes (mph)
	WindDirAvgLast2Min     *float64 `json:"wind_dir_scalar_avg_last_2_min"`   // WindDirAvgLast2Min is average wind direction over last 2 minutes (°)
	WindSpeedHighLast2Min  *float64 `json:"wind_speed_hi_last_2_min"`         // WindSpeedHighLast2Min is max gust over last 2 minutes (mph)
	WindDirAtHighLast2Min  *float64 `json:"wind_dir_at_hi_speed_last_2_min"`  // WindDirAtHighLast2Min is max gust direction over last 2 minutes (°)
	WindSpeedAvgLast10Min  *float64 `json:"wind_speed_avg_last_10_min"`       // WindSpeedAvgLast10Min is average wind over last 10 minutes (mph)
	WindDirAvgLast10Min    *float64 `json:"wind_dir_scalar_avg_last_10_min"`  // WindDirAvgLast10Min is average wind dir over last 10 minutes (°)
	WindSpeedHighLast10Min *float64 `json:"wind_speed_hi_last_10_min"`        // WindSpeedHighLast10Min is max gust over last 10 minutes (mph)
	WindDirAtHighLast10Min *float64 `json:"wind_dir_at_hi_speed_last_10_min"` // WindDirAtHighLast10Min is max gust direction over last 10 minutes (°)

	RainSize              *float64 `json:"rain_size"`                // RainSize is size of rain collector (1: 0.01", 2: 0.2mm)
	RainRateLast          *float64 `json:"rain_rate_last"`           // RainRateLast is most recent rain rate (count/hour)
	RainRateHigh          *float64 `json:"rain_rate_hi"`             // RainRateHigh is highest rain rate over last minute (count/hour)
	RainLast15Min         *float64 `json:"rainfall_last_15_min"`     // RainLast15Min is rain count in last 15 minutes (count)
	RainRateHighLast15Min *float64 `json:"rain_rate_hi_last_15_min"` // RainRateHighLast15Min is highest rain count rate over last 15 minutes (count/hour)
	RainLast60Min         *float64 `json:"rainfall_last_60_min"`     // RainLast60Min is rain count over last 60 minutes (count)
	RainLast24Hour        *float64 `json:"rainfall_last_24_hr"`      // RainLast24Hour is rain count over last 24 hours (count)
	RainStorm             *float64 `json:"rain_storm"`               // RainStorm is rain since last 24 hour break in rain (count)
	RainStormStartAt      *Time    `json:"rain_storm_start_at"`      // RainStormStartAt is time of rain storm start

	SolarRad *float64 `json:"solar_rad"` // SolarRad is solar radiation (W/m²)
	UVIndex  *float64 `json:"uv_index"`  // UVIndex is solar UV index

	RXState          *SignalState  `json:"rx_state"`           // RXState is ISS receiver status
	TransBatteryFlag *BatteryState `json:"trans_battery_flag"` // TransBatteryFlag is ISS battery status

	RainfallDaily        *float64 `json:"rainfall_daily"`           // RainfallDaily is total rain since midnight (count)
	RainfallMonthly      *float64 `json:"rainfall_monthly"`         // RainfallMonthly is total rain since first of month (count)
	RainfallYear         *float64 `json:"rainfall_year"`            // RainfallYear is total rain since first of year (count)
	RainStormLast        *float64 `json:"rain_storm_last"`          // RainStormLast is rain since last 24 hour break in rain (count)
	RainStormLastStartAt *Time    `json:"rain_storm_last_start_at"` // RainStormLastStartAt is time of last rain storm start
	RainStormLastEndAt   *Time    `json:"rain_storm_last_end_at"`   // rainStormLastEndAt is time of last rain storm end
}

// WeatherLSSBarometer is weather conditions from LSS barometer sensor.
type WeatherLSSBarometer struct {
	BarometerSeaLevel *float64 `json:"bar_sea_level"` // BarometerSeaLevel is barometer reading with elevation adjustment (inches)
	BarometerTrend    *float64 `json:"bar_trend"`     // BarometerTrend is 3 hour barometric trend (inches)
	BarometerAbsolute *float64 `json:"bar_absolute"`  // BarometerAbsolute is barometer reading at current elevation (inches)
}

// WeatherLSSTempRh is weather conditions from LSS temperature and humidity
// sensors.
type WeatherLSSTempRh struct {
	TemperatureIndoor *float64 `json:"temp_in"`       // TemperatureIndoor is indoor temp (°F)
	HumidityIndoor    *float64 `json:"hum_in"`        // HumidityIndoor is indoor humidity (%)
	DewPointIndoor    *float64 `json:"dew_point_in"`  // DewPointIndoor is indoor dewpoint (°F)
	HeatIndexIndoor   *float64 `json:"heat_index_in"` // HeatIndexIndoor is indoor heat index (°F)
}

// ParseHTTP parses weather conditions retrieved over HTTP. It returns an error
// if the payload format is not valid.
func ParseHTTP(data []byte) (*ConditionsHTTP, error) {
	var out ConditionsHTTP
	err := json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// UnmarshalJSON unmarshals the weather condition entry retrieved over HTTP. It
// returns an error if the payload format is not valid.
func (e *EntryHTTP) UnmarshalJSON(payload []byte) error {
	// read partial condition header
	var c conditionHeader
	err := json.Unmarshal(payload, &c)
	if err != nil {
		return err
	}
	// place retrieved values in main body
	e.LogicalSensorID = c.LogicalSensorID
	e.DataStructureType = c.DataStructureType

	// parse remaining using specified data structure
	switch *c.DataStructureType {
	case RecordISS:
		e.Values = &WeatherISS{}
	case RecordLSSBarometer:
		e.Values = &WeatherLSSBarometer{}
	case RecordLSSTempRh:
		e.Values = &WeatherLSSTempRh{}
	default:
		e.Values = nil
	}

	// recursively parse inner values
	return json.Unmarshal(payload, e.Values)
}

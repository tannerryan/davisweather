// Copyright (c) 2020 Tanner Ryan. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

package parser

import (
	"encoding/json"
)

// ConditionsUDP is for parsing weather conditions received over UDP.
type ConditionsUDP struct {
	DeviceID   string     `json:"did"`        // DeviceID is unique device ID
	Timestamp  *Time      `json:"ts"`         // Timestamp is device timestamp (epoch seconds)
	Conditions []EntryUDP `json:"conditions"` // Conditions is a list of conditions
}

// EntryUDP is for parsing UDP weather data.
type EntryUDP struct {
	LogicalSensorID        *int        `json:"lsid"`                             // LogicalSensorID is logical sensor ID
	DataStructureType      *RecordType `json:"data_structure_type"`              // DataStructureType indicates data structure
	TransmitterID          *int        `json:"txid"`                             // TransmitterID is ID of transmitter
	WindSpeedLast          *float64    `json:"wind_speed_last"`                  // WindSpeedLast is most recent wind speed (mph)
	WindDirLast            *float64    `json:"wind_dir_last"`                    // WindDirLast is most recent wind direction (°)
	RainSize               *float64    `json:"rain_size"`                        // RainSize is size of rain collector (1: 0.01", 2: 0.2mm)
	RainRateLast           *float64    `json:"rain_rate_last"`                   // RainRateLast is most recent rain rate (count/hour)
	RainLast15Min          *float64    `json:"rain_15_min"`                      // RainLast15Min is rain count over last 15 minutes (count)
	RainLast60Min          *float64    `json:"rain_60_min"`                      // RainLast60Min is rain count over last 60 minutes (count)
	RainLast24Hour         *float64    `json:"rain_24_hr"`                       // RainLast24Hour is rain count over last 24 hours (count)
	RainStorm              *float64    `json:"rain_storm"`                       // RainStorm is rain since last 24 hour break in rain (count)
	RainStormStartAt       *Time       `json:"rain_storm_start_at"`              // RainStormStartAt is time of rainstorm start
	RainfallDaily          *float64    `json:"rainfall_daily"`                   // RainfallDaily is total rain since midnight (count)
	RainfallMonthly        *float64    `json:"rainfall_monthly"`                 // RainfallMonthly is total rain since first of month (count)
	RainfallYear           *float64    `json:"rainfall_year"`                    // RainfallYear is total rain since first of year (count)
	WindSpeedHighLast10Min *float64    `json:"wind_speed_hi_last_10_min"`        // WindSpeedHighLast10Min is max gust over last 10 minutes (mph)
	WindDirAtHighLast10Min *float64    `json:"wind_dir_at_hi_speed_last_10_min"` // WindDirAtHighLast10Min is max gust direction over last 10 minutes (°)
}

// ParseUDP parses weather conditions received over UDP. It returns an error if
// the payload format is not valid.
func ParseUDP(payload []byte) (*ConditionsUDP, error) {
	var out ConditionsUDP
	err := json.Unmarshal(payload, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

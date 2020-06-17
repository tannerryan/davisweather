// Copyright (c) 2020 Tanner Ryan. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

package davisweather

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/thetannerryan/davisweather/parser"
)

// Report is the latest weather report.
type Report struct {
	DeviceID  string    `json:"deviceID"`  // DeviceID is unique device ID
	Timestamp time.Time `json:"timestamp"` // Timestamp is the time the Report was last modified

	Temperature *float64 `json:"temperature"` // Temperature (°F)
	Humidity    *float64 `json:"humidity"`    // Humidity (%RH)
	Dewpoint    *float64 `json:"dewpoint"`    // Dewpoint (°F)
	Wetbulb     *float64 `json:"wetbulb"`     // Wetbulb (°F)
	HeatIndex   *float64 `json:"heatindex"`   // HeatIndex (°F)
	WindChill   *float64 `json:"windchill"`   // WindChill (°F)
	THWIndex    *float64 `json:"thwIndex"`    // THWIndex is "feels like" (°F)
	THSWIndex   *float64 `json:"thswIndex"`   // THWSIndex is "feels like" including solar (°F)

	WindSpeedLast          *float64 `json:"windSpeedLast"`          // WindSpeedLast is most recent wind speed (mph)
	WindDirLast            *float64 `json:"windDirLast"`            // WindDirLast is most recent wind direction (°)
	WindSpeedAvgLast1Min   *float64 `json:"windSpeedAvg1Min"`       // WindSpeedAvgLast1Min is average wind over last minute (mph)
	WindDirAvgLast1Min     *float64 `json:"windDirAvg1Min"`         // WindDirAvgLast1Min is average wind direction over last minute (°)
	WindSpeedAvgLast2Min   *float64 `json:"windSpeedAvg2Min"`       // WindSpeedAvgLast2Min is average wind over last 2 minutes (mph)
	WindDirAvgLast2Min     *float64 `json:"windDirAvg2Min"`         // WindDirAvgLast2Min is average wind direction over last 2 minutes (°)
	WindSpeedHighLast2Min  *float64 `json:"windGustSpeedLast2Min"`  // WindSpeedHighLast2Min is max gust over last 2 minutes (mph)
	WindDirAtHighLast2Min  *float64 `json:"windGustDirLast2Min"`    // WindDirAtHighLast2Min is max gust direction over last 2 minutes (°)
	WindSpeedAvgLast10Min  *float64 `json:"windSpeedAvg10Min"`      // WindSpeedAvgLast10Min is average wind over last 10 minutes (mph)
	WindDirAvgLast10Min    *float64 `json:"windDirAvg10Min"`        // WindDirAvgLast10Min is average wind dir over last 10 minutes (°)
	WindSpeedHighLast10Min *float64 `json:"windGustSpeedLast10Min"` // WindSpeedHighLast10Min is max gust over last 10 minutes (mph)
	WindDirAtHighLast10Min *float64 `json:"windGustDirLast10Min"`   // WindDirAtHighLast10Min is max gust direction over last 10 minutes (°)

	RainSize              *float64   `json:"rainSize"`              // RainSize is size of rain collector (1: 0.01", 2: 0.2mm)
	RainRateLast          *float64   `json:"rainRateLast"`          // RainRateLast is most recent rain rate (count/hour)
	RainRateHigh          *float64   `json:"rainRateHigh"`          // RainRateHigh is highest rain rate over last minute (count/hour)
	RainLast15Min         *float64   `json:"rainLast15Min"`         // RainLast15Min is rain count in last 15 minutes (count)
	RainRateHighLast15Min *float64   `json:"rainRateHighLast15Min"` // RainRateHighLast15Min is highest rain count rate over last 15 minutes (count/hour)
	RainLast60Min         *float64   `json:"rainLast60Min"`         // RainLast60Min is rain count over last 60 minutes (count)
	RainLast24Hour        *float64   `json:"rainLast24Hour"`        // RainLast24Hour is rain count over last 24 hours (count)
	RainStorm             *float64   `json:"rainStorm"`             // RainStorm is rain since last 24 hour break in rain (count)
	RainStormStartAt      *time.Time `json:"rainStormStart"`        // RainStormStartAt is time of rain storm start

	SolarRad *float64 `json:"solarRad"` // SolarRad is solar radiation (W/m²)
	UVIndex  *float64 `json:"uvIndex"`  // UVIndex is solar UV index

	RXState          string `json:"signal"`  // RXState is ISS receiver status
	TransBatteryFlag string `json:"battery"` // TransBatteryFlag is ISS battery status

	RainfallDaily        *float64   `json:"rainDaily"`          // RainfallDaily is total rain since midnight (count)
	RainfallMonthly      *float64   `json:"rainMonthly"`        // RainfallMonthly is total rain since first of month (count)
	RainfallYear         *float64   `json:"rainYear"`           // RainfallYear is total rain since first of year (count)
	RainStormLast        *float64   `json:"rainStormLast"`      // RainStormLast is rain since last 24 hour break in rain (count)
	RainStormLastStartAt *time.Time `json:"rainStormLastStart"` // RainStormLastStartAt is time of last rain storm start
	RainStormLastEndAt   *time.Time `json:"rainStormLastEnd"`   // rainStormLastEndAt is time of last rain storm end

	BarometerSeaLevel *float64 `json:"barometerSeaLevel"` // BarometerSeaLevel is barometer reading with elevation adjustment (inches)
	BarometerTrend    *float64 `json:"barometerTrend"`    // BarometerTrend is 3 hour barometric trend (inches)
	BarometerAbsolute *float64 `json:"barometerAbsolute"` // BarometerAbsolute is barometer reading at current elevation (inches)

	TemperatureIndoor *float64 `json:"indoorTemperature"` // TemperatureIndoor is indoor temp (°F)
	HumidityIndoor    *float64 `json:"indoorHumidity"`    // HumidityIndoor is indoor humidity (%)
	DewPointIndoor    *float64 `json:"indoorDewpoint"`    // DewPointIndoor is indoor dewpoint (°F)
	HeatIndexIndoor   *float64 `json:"indoorHeatIndex"`   // HeatIndexIndoor is indoor heat index (°F)

	notify       chan bool   // notify emits a boolean when the Report contents are modified
	verbose      bool        // verbose enables Report logging to stdout
	lastChecksum string      // lastChecksum is MD5 checksum of the Report state
	lastBytes    []byte      // lastBytes is the JSON representation of the Report state
	mutex        *sync.Mutex // mutex is for atomic report actions
}

// NewReport returns a new Report state and a notification channel. The channel
// emits a bool when the Report contents have been modified. Verbose enables
// Report logging.
func NewReport(verbose bool) (*Report, chan bool) {
	r := &Report{
		notify:       make(chan bool, 1),
		verbose:      verbose,
		lastChecksum: "",
		lastBytes:    nil,
		mutex:        &sync.Mutex{},
	}
	return r, r.notify
}

// UpdateHTTP atomically updates the Report state using weather conditions
// retrieved over HTTP. It returns an error if the atomic update action fails.
func (r *Report) UpdateHTTP(new *parser.ConditionsHTTP) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// ensure no error in conditions
	if new.Error != nil {
		return errors.New(new.Error.Message)
	}

	// set report header
	r.DeviceID = new.Data.DeviceID

	// iterate over all provided conditions, load conditions into report
	for _, c := range new.Data.Conditions {
		structure := *c.DataStructureType

		switch structure {
		case parser.RecordISS:
			v := c.Values.(*parser.WeatherISS)
			r.processISS(v)
		case parser.RecordLSSBarometer:
			v := c.Values.(*parser.WeatherLSSBarometer)
			r.processLSSBarometer(v)
		case parser.RecordLSSTempRh:
			v := c.Values.(*parser.WeatherLSSTempRh)
			r.processLSSTempRh(v)
		}
	}

	return r.updateHook(parser.UpdateHTTP, new.Data.Timestamp.Time())
}

// UpdateUDP atomically updates the Report state using weather conditions
// received over UDP. It returns an error if the atomic update action fails.
func (r *Report) UpdateUDP(new *parser.ConditionsUDP) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// set report header
	r.DeviceID = new.DeviceID

	// iterate over all conditions, load conditions into report
	for _, c := range new.Conditions {
		r.WindSpeedLast = c.WindSpeedLast
		r.WindDirLast = c.WindDirLast

		r.RainSize = c.RainSize
		r.RainRateLast = c.RainRateLast
		r.RainLast15Min = c.RainLast15Min
		r.RainLast60Min = c.RainLast60Min
		r.RainLast24Hour = c.RainLast24Hour
		r.RainStorm = c.RainStorm
		if c.RainStormStartAt != nil {
			t := c.RainStormStartAt.Time()
			r.RainStormStartAt = &t
		}
		r.RainfallDaily = c.RainfallDaily
		r.RainfallMonthly = c.RainfallMonthly
		r.RainfallYear = c.RainfallYear

		r.WindSpeedHighLast10Min = c.WindSpeedHighLast10Min
		r.WindDirAtHighLast10Min = c.WindDirAtHighLast10Min
	}

	return r.updateHook(parser.UpdateUDP, new.Timestamp.Time())
}

// UpdateJSON atomically updates the Report state using a provided JSON payload.
// It returns an error if the JSON does not match the Report structure or if the
// atomic update action fails.
func (r *Report) UpdateJSON(payload []byte) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// attempt parsing
	var n Report
	err := json.Unmarshal(payload, &n)
	if err != nil {
		return err
	}

	// synchronize all data fields
	r.DeviceID = n.DeviceID

	r.Temperature = n.Temperature
	r.Humidity = n.Humidity
	r.Dewpoint = n.Dewpoint
	r.Wetbulb = n.Wetbulb
	r.HeatIndex = n.HeatIndex
	r.WindChill = n.WindChill
	r.THWIndex = n.THWIndex
	r.THSWIndex = n.THSWIndex

	r.WindSpeedLast = n.WindSpeedLast
	r.WindDirLast = n.WindDirLast
	r.WindSpeedAvgLast1Min = n.WindSpeedAvgLast1Min
	r.WindDirAvgLast1Min = n.WindDirAvgLast1Min
	r.WindSpeedAvgLast2Min = n.WindSpeedAvgLast2Min
	r.WindDirAvgLast2Min = n.WindDirAvgLast2Min
	r.WindSpeedHighLast2Min = n.WindSpeedHighLast2Min
	r.WindDirAtHighLast2Min = n.WindDirAtHighLast2Min
	r.WindSpeedAvgLast10Min = n.WindSpeedAvgLast10Min
	r.WindDirAvgLast10Min = n.WindDirAvgLast10Min
	r.WindSpeedHighLast10Min = n.WindSpeedHighLast10Min
	r.WindDirAtHighLast10Min = n.WindDirAtHighLast10Min

	r.RainSize = n.RainSize
	r.RainRateLast = n.RainRateLast
	r.RainRateHigh = n.RainRateHigh
	r.RainLast15Min = n.RainLast15Min
	r.RainRateHighLast15Min = n.RainRateHighLast15Min
	r.RainLast60Min = n.RainLast60Min
	r.RainLast24Hour = n.RainLast24Hour
	r.RainStorm = n.RainStorm
	r.RainStormStartAt = n.RainStormStartAt

	r.SolarRad = n.SolarRad
	r.UVIndex = n.UVIndex

	r.RXState = n.RXState
	r.TransBatteryFlag = n.TransBatteryFlag

	r.RainfallDaily = n.RainfallDaily
	r.RainfallMonthly = n.RainfallMonthly
	r.RainfallYear = n.RainfallYear
	r.RainStormLast = n.RainStormLast
	r.RainStormLastStartAt = n.RainStormLastStartAt
	r.RainStormLastEndAt = n.RainStormLastEndAt

	r.BarometerSeaLevel = n.BarometerSeaLevel
	r.BarometerTrend = n.BarometerTrend
	r.BarometerAbsolute = n.BarometerAbsolute

	r.TemperatureIndoor = n.TemperatureIndoor
	r.HumidityIndoor = n.HumidityIndoor
	r.DewPointIndoor = n.DewPointIndoor
	r.HeatIndexIndoor = n.HeatIndexIndoor

	return r.updateHook(parser.UpdateJSON, time.Now())
}

// Copy generates and returns a new deep copy of the Report state. It returns an
// error if the copy operation fails.
func (r *Report) Copy() (*Report, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// marshal current report
	buff, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	report, notify := NewReport(r.verbose)
	// unmarshal buffer into report
	err = json.Unmarshal(buff, &report)
	if err != nil {
		return nil, err
	}
	// sync private variables
	report.notify = notify
	report.lastChecksum = r.lastChecksum
	report.lastBytes = r.lastBytes
	report.mutex = &sync.Mutex{}

	return report, nil
}

// JSON returns the JSON representation of the Report state and the Report was
// last updated.
func (r *Report) JSON() []byte {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.lastBytes
}

// Encode returns a zlib encoded weather report.
func (r *Report) Encode() []byte {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// stream last bytes into compressed buffer
	var buff bytes.Buffer
	stream := zlib.NewWriter(&buff)
	stream.Write(r.lastBytes)
	stream.Close()

	// return compressed bytes
	return buff.Bytes()
}

// Decode updates the Report using a zlib encoded weather report. It returns an
// error if the provided payload is not valid.
func (r *Report) Decode(payload []byte) error {
	// uncompress data
	buff := bytes.NewReader(payload)
	stream, err := zlib.NewReader(buff)
	if err != nil {
		return err
	}
	defer stream.Close()

	// read uncompressed data
	report, err := ioutil.ReadAll(stream)
	if err != nil {
		return err
	}
	return r.UpdateJSON(report)
}

// updateHook is called after UpdateHTTP, UpdateUDP, and UpdateJSON. If the
// Report state contents have changed, it updates the lastBytes, lastChecksum,
// and Timestamp fields, emitting a bool on the notify channel if the channel is
// not full. If the Report state contents are not changed, nothing happens. It
// returns an error if the Report state checksum fails.
func (r *Report) updateHook(method parser.UpdateMethod, timestamp time.Time) error {
	// calculate checksum of latest report
	newChecksum, _, err := r.checksum()
	if err != nil {
		return err
	}
	// ignore if no ISS transmitter or ISS transmitter battery flag is given
	if r.RXState == "" || r.TransBatteryFlag == "" {
		return nil
	}
	// only update internals, timestamp, and notify if new content
	if newChecksum != r.lastChecksum {
		// update last timestamp, last bytes and checksum
		r.Timestamp = timestamp
		r.lastChecksum, r.lastBytes, _ = r.checksum()
		select {
		case r.notify <- true: // attempt to notify
			if r.verbose {
				log.Println("[davisweather report] new data from", r.DeviceID, method)
			}
		default: // already notified on channel
			if r.verbose {
				log.Println("[davisweather report] new data from", r.DeviceID, method, "(downstream pressure on Notify)")
			}
		}
		return nil
	}
	if r.verbose {
		log.Println("[davisweather report] no new data from", r.DeviceID, method)
	}
	return nil
}

// checksum return an MD5 checksum of the Report state and the JSON
// representation of the Report state. It returns an error if the checksum
// fails.
func (r *Report) checksum() (string, []byte, error) {
	// marshal current state
	buff, err := json.Marshal(r)
	if err != nil {
		return "", nil, err
	}
	// perform sum
	hash := md5.Sum(buff)
	return hex.EncodeToString(hash[:]), buff, nil
}

// processISS synchronizes the Report state with the provided ISS weather
// conditions.
func (r *Report) processISS(v *parser.WeatherISS) {
	r.Temperature = v.Temperature
	r.Humidity = v.Humidity
	r.Dewpoint = v.Dewpoint
	r.Wetbulb = v.Wetbulb
	r.HeatIndex = v.HeatIndex
	r.WindChill = v.WindChill
	r.THWIndex = v.THWIndex
	r.THSWIndex = v.THSWIndex

	r.WindSpeedLast = v.WindSpeedLast
	r.WindDirLast = v.WindDirLast
	r.WindSpeedAvgLast1Min = v.WindSpeedAvgLast1Min
	r.WindDirAvgLast1Min = v.WindDirAvgLast1Min
	r.WindSpeedAvgLast2Min = v.WindSpeedAvgLast2Min
	r.WindDirAvgLast2Min = v.WindDirAvgLast2Min
	r.WindSpeedHighLast2Min = v.WindSpeedHighLast2Min
	r.WindDirAtHighLast2Min = v.WindDirAtHighLast2Min
	r.WindSpeedAvgLast10Min = v.WindSpeedAvgLast10Min
	r.WindDirAvgLast10Min = v.WindDirAvgLast10Min
	r.WindSpeedHighLast10Min = v.WindSpeedHighLast10Min
	r.WindDirAtHighLast10Min = v.WindDirAtHighLast10Min

	r.RainSize = v.RainSize
	r.RainRateLast = v.RainRateLast
	r.RainRateHigh = v.RainRateHigh
	r.RainLast15Min = v.RainLast15Min
	r.RainRateHighLast15Min = v.RainRateHighLast15Min
	r.RainLast60Min = v.RainLast60Min
	r.RainLast24Hour = v.RainLast24Hour
	r.RainStorm = v.RainStorm
	if v.RainStormStartAt != nil {
		t := v.RainStormStartAt.Time()
		r.RainStormStartAt = &t
	}

	r.SolarRad = v.SolarRad
	r.UVIndex = v.UVIndex

	if v.RXState != nil {
		switch *v.RXState {
		case parser.SignalSynced:
			r.RXState = "Synced"
		case parser.SignalRescan:
			r.RXState = "Rescan"
		case parser.SignalLost:
			r.RXState = "Lost"
		}
	}
	if v.TransBatteryFlag != nil {
		switch *v.TransBatteryFlag {
		case parser.BatteryNominal:
			r.TransBatteryFlag = "Nominal"
		case parser.BatteryWarning:
			r.TransBatteryFlag = "Warning"
		}
	}

	r.RainfallDaily = v.RainfallDaily
	r.RainfallMonthly = v.RainfallMonthly
	r.RainfallYear = v.RainfallYear
	r.RainStormLast = v.RainStormLast
	if v.RainStormLastStartAt != nil {
		t := v.RainStormLastStartAt.Time()
		r.RainStormLastStartAt = &t
	}
	if v.RainStormLastEndAt != nil {
		t := v.RainStormLastEndAt.Time()
		r.RainStormLastEndAt = &t
	}
}

// processLSSBarometer synchronizes the Report state with the provided LSS
// barometer weather conditions.
func (r *Report) processLSSBarometer(v *parser.WeatherLSSBarometer) {
	r.BarometerSeaLevel = v.BarometerSeaLevel
	r.BarometerTrend = v.BarometerTrend
	r.BarometerAbsolute = v.BarometerAbsolute
}

// processLSSBarometer synchronizes the Report state with the provided LSS
// temperature humidity weather conditions.
func (r *Report) processLSSTempRh(v *parser.WeatherLSSTempRh) {
	r.TemperatureIndoor = v.TemperatureIndoor
	r.HumidityIndoor = v.HumidityIndoor
	r.DewPointIndoor = v.DewPointIndoor
	r.HeatIndexIndoor = v.HeatIndexIndoor
}

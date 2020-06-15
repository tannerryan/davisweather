// Copyright (c) 2020 Tanner Ryan. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

package parser

// RecordType indicates the weather condition type returned from HTTP.
type RecordType int

const (
	// RecordISS is weather conditions from integrated sensor suite
	RecordISS RecordType = 1
	// RecordLSSBarometer is weather conditions from LSS barometer
	RecordLSSBarometer RecordType = 3
	// RecordLSSTempRh is weather condition from LSS temperature + humidity
	RecordLSSTempRh RecordType = 4
)

// SignalState indicates the WLL receiver status.
type SignalState int

const (
	// SignalSynced means WLL unit is locked onto ISS signal
	SignalSynced SignalState = 0
	// SignalRescan means WLL unit is scanning for ISS signal
	SignalRescan SignalState = 1
	// SignalLost means WLL unit failed to locate ISS signal
	SignalLost SignalState = 2
)

// BatteryState indicates the ISS battery status.
type BatteryState int

const (
	// BatteryNominal means ISS battery does not require replacement
	BatteryNominal BatteryState = 0
	// BatteryWarning means ISS battery requires replacement
	BatteryWarning BatteryState = 1
)

// UpdateMethod indicates how the Report state is updated
type UpdateMethod string

const (
	// UpdateHTTP is updating Report state with data retrieved from HTTP
	UpdateHTTP UpdateMethod = "http"
	// UpdateUDP is updating Report state with data received from UDP
	UpdateUDP UpdateMethod = "udp"
	// UpdateJSON is updating Report state with a JSON payload
	UpdateJSON UpdateMethod = "json"
)

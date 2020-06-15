// Copyright (c) 2020 Tanner Ryan. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

package parser

import "encoding/json"

// BroadcastResponse is for parsing the response of enabling UDP broadcasts.
type BroadcastResponse struct {
	ConnInfo *ConnInfo `json:"data"`  // ConnInfo contains the UDP broadcast parameters
	Error    *Error    `json:"error"` // Error is a general error message
}

// ConnInfo are the UDP broadcast parameters.
type ConnInfo struct {
	Port     int `json:"broadcast_port"` // Port is the UDP port of weather broadcasts
	Duration int `json:"duration"`       // Duration is how long UDP broadcast is enabled for
}

// ParseBroadcastResponse parses the response after requesting UDP weather
// broadcasts. It returns an error if the payload format is not valid.
func ParseBroadcastResponse(payload []byte) (*BroadcastResponse, error) {
	var out BroadcastResponse
	err := json.Unmarshal(payload, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

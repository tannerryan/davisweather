// Copyright (c) 2020 Tanner Ryan. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

package parser

// Error is a generic error message provided by WLL unit.
type Error struct {
	Code    int    `json:"code"`    // Code is the error code
	Message string `json:"message"` // Message is an error message
}

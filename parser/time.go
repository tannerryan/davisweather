// Copyright (c) 2020 Tanner Ryan. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

package parser

import (
	"strconv"
	"time"
)

// Time is an alias for parsing epoch timestamps.
type Time time.Time

// UnmarshalJSON is to unmarshal epoch timestamp into time.Time. It returns an
// error if the provided timestamp is invalid.
func (t *Time) UnmarshalJSON(payload []byte) error {
	val, err := strconv.ParseInt(string(payload), 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(val, 0)
	return nil
}

// Time returns the underlining time.Time instance.
func (t Time) Time() time.Time {
	return time.Time(t)
}

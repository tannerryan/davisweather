// Copyright (c) 2020 Tanner Ryan. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

package main

import (
	"context"
	"log"

	"github.com/thetannerryan/davisweather"
)

func main() {
	ctx := context.Background()
	client := davisweather.Managed(ctx, true)

	for {
		<-client.Notify
		report, err := client.Report()
		if err != nil {
			panic(err)
		}
		log.Println(*report.Temperature)
	}
}

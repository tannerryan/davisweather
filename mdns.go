// Copyright (c) 2020 Tanner Ryan. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

package davisweather

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
)

const (
	// mDNSInstance is mDNS instance name
	mDNSInstance = "_weatherlinklive"
	// mDNSService is mDNS service name
	mDNSService = "_tcp."
	// mDNSDomain is mDNS domain name
	mDNSDomain = "local."
	// mDNSTimeout is the mDNS discovery timeout
	mDNSTimeout = 15 * time.Second
)

var (
	// mDNSInterval is sleep between mDNS discovery (gets modified to TTL)
	mDNSInterval = 5 * time.Second
)

// wllUnit represents the network configuration of a WLL unit
type wllUnit zeroconf.ServiceEntry

// GetURL generates an HTTP URL from the wllUnit, or returns an empty string.
func (u *wllUnit) GetURL() string {
	if len(u.AddrIPv6) > 0 {
		return fmt.Sprintf("http://[%s]:%d", u.AddrIPv6[0].String(), u.Port)
	}
	if len(u.AddrIPv4) > 0 {
		return fmt.Sprintf("http://%s:%d", u.AddrIPv4[0].String(), u.Port)
	}
	if u.HostName != "" {
		return fmt.Sprintf("http://%s:%d", u.HostName, u.Port)
	}
	return ""
}

// discovery runs the mDNS discover routine on set intervals to locate the WLL
// unit. The resolved cancel function is called when a WLL unit is found.
func (c *Client) discovery(ctx context.Context, resolved context.CancelFunc) {
	// goroutine monitoring
	defer c.wg.Done()

	// prevent leaking of context if loop terminates
	defer resolved()

	for {
		c.println("[davisweather mdns] performing autodiscovery of WeatherLink Live unit")
		recover := mDNSInterval / 2
		err := c.mDNSDiscover(ctx, resolved)
		if err != nil {
			c.println("[davisweather mdns] failed to perform autodiscovery, retrying in", recover)
			time.Sleep(recover)
			continue
		}
		// sleep or terminate
		select {
		case <-ctx.Done():
			// Client was terminated
			c.println("[davisweather mdns] terminating event loop")
			return
		case <-time.After(mDNSInterval):
			// sleep for next iteration
		}
	}
}

// mDNSDiscover is called on regular intervals to discover the WLL unit. It
// returns an error if the mDNS discover process fails. The resolved cancel
// function is called when a WLL unit is found.
func (c *Client) mDNSDiscover(ctx context.Context, resolved context.CancelFunc) error {
	// initialize UDP resolver
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return err
	}
	// terminate discovery when device is found or timeout
	deadline := time.Now().Add(mDNSTimeout)
	ctx, done := context.WithDeadline(ctx, deadline)
	defer done()

	// start consume responder channel
	responders := make(chan *zeroconf.ServiceEntry)
	go mDNSLoop(c, responders, done, resolved)

	// perform mDNS lookup (closes responders channel on ctx Done)
	err = resolver.Lookup(ctx, mDNSInstance, mDNSService, mDNSDomain, responders)
	if err != nil {
		return err
	}
	// wait until device is found or timeout
	<-ctx.Done()

	return nil
}

// mDNSLoop is called by the mDNSDiscover process. It loops over the
// ServiceEntry channel until the WLL unit is located. If the unit is located,
// the unit in Client is updated and the mDNSInterval is updated to mDNS TTL.
func mDNSLoop(c *Client, responders <-chan *zeroconf.ServiceEntry, done context.CancelFunc, resolved context.CancelFunc) {
	start := time.Now()
	for r := range responders {
		if strings.Contains(r.ServiceRecord.Instance, mDNSInstance) {
			// calculate discovery duration
			duration := time.Now().Sub(start)

			// generate wllUnit and update Client
			u := wllUnit(*r)
			c.unit = &u

			// update mDNS interval to TTL
			mDNSInterval = time.Duration(r.TTL) * time.Second

			// location printing
			if len(u.AddrIPv6) > 0 {
				c.printf("[davisweather mdns] found WeatherLink Live unit in %.2fs at %s:%d",
					duration.Seconds(), u.AddrIPv6[0].String(), u.Port)
			} else if len(u.AddrIPv4) > 0 {
				c.printf("[davisweather mdns] found WeatherLink Live unit in %.2fs at %s:%d",
					duration.Seconds(), u.AddrIPv4[0].String(), u.Port)
			}

			c.println("[davisweather mdns] reperforming autodiscovery in", mDNSInterval)

			// notify caller process is done
			done()
			resolved()
			return
		}
	}
}

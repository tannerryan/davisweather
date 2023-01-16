// Copyright (c) 2020 Tanner Ryan. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

package davisweather

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/tannerryan/davisweather/parser"
)

const (
	// engineIntervalHTTP is how often to poll for HTTP weather conditions
	engineIntervalHTTP = 10300 * time.Millisecond // 10.3 seconds
	// engineIntervalWatchdog is how often to run the UDP watchdog
	engineIntervalWatchdog = udpDeadline / 2.0

	// httpTimeout is the HTTP client timeout
	httpTimeout = 3 * time.Second
	// routeConditions is route for fetching weather conditions
	routeConditions = "/v1/current_conditions"
	// routeBroadcastResponse is route for fetching broadcast response
	routeBroadcastResponse = "/v1/real_time"

	// udpDeadline is UDP read deadline before sending broadcast response
	udpDeadline = 15 * time.Second
	// udpDuration is how long UDP broadcasts are enabled for
	udpDuration = 4 * time.Hour
	// udpBufferSize is buffer size for reading UDP messages
	udpBufferSize = 2048
)

// engine starts the UDP and HTTP event loops after the mDNS autodiscovery is
// completed. If the unit in Client is defined (unmanaged client), the engine
// starts the UDP and HTTP loops without delay.
func (c *Client) engine(ctx context.Context, mDNS context.Context) {
	// goroutine monitoring
	defer c.wg.Done()

	// stall engine until mDNS has resolved
	if c.unit == nil && mDNS != nil {
		c.println("[davisweather] waiting for mDNS autodiscovery")
		<-mDNS.Done()
	}
	select {
	case <-ctx.Done():
		// received termination signal
		return
	default:
		// start HTTP and UDP event loops
		c.println("[davisweather] initializing UDP and HTTP event loops")
		c.wg.Add(2)
		go c.httpEventLoop(ctx)
		go c.udpEventLoop(ctx)

		// received termination signal
		<-ctx.Done()
		return
	}
}

// httpEventLoop perodically retrieves weather conditions over HTTP and updates
// the Report state in Client.
func (c *Client) httpEventLoop(ctx context.Context) {
	// goroutine monitoring
	defer c.wg.Done()

	// WLL does can't support concurrent HTTP, allow UDP to get first request
	time.Sleep(5 * time.Second)

	// initialize event timer
	eventTimer := time.NewTimer(engineIntervalHTTP)
	defer eventTimer.Stop()
	c.println("[davisweather http] initializing, fetching weather conditions every", engineIntervalHTTP)

	for {
		eventTimer.Reset(engineIntervalHTTP)

		// fetch latest conditions
		conditions, err := c.fetchConditionsHTTP(ctx)
		if err != nil {
			c.println("[davisweather http] failed to fetch conditions", err)
		} else {
			// update Report state
			err = c.report.UpdateHTTP(conditions)
			if err != nil {
				c.println("[davisweather http] failed to update Report", err)
			}
		}
		// terminate or sleep
		select {
		case <-ctx.Done():
			c.println("[davisweather http] terminating event loop")
			eventTimer.Stop()
			return
		case <-eventTimer.C:
		}
	}
}

// udpEventLoop receives weather conditions over UDP and updates the Report
// state in Client.
func (c *Client) udpEventLoop(ctx context.Context) {
	// goroutine monitoring
	defer c.wg.Done()

	// UDP context to notify event loop UDP port is acquired
	udpCtx, udpDone := context.WithCancel(ctx)
	defer udpDone()

	// start UDP watchdog
	go c.udpWatchdog(ctx, udpDone)

	for {
		// check if UDP port provided
		if c.udpPort == 0 {
			// terminate or wait for UDP context to resolve
			select {
			case <-ctx.Done():
				c.println("[davisweather udp] terminating event loop")
				return
			case <-udpCtx.Done():
			}
		}
		// terminate or wait error delay
		select {
		case <-ctx.Done():
			c.println("[davisweather udp] terminating event loop")
			return
		case <-time.After(engineIntervalWatchdog):
		}

		// generate connection address
		connAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(c.udpPort))
		if err != nil {
			c.println("[davisweather udp] failed to generate UDP address, retrying in", engineIntervalWatchdog)
			continue
		}
		// establish connection to UDP socket
		conn, err := net.ListenUDP("udp", connAddr)
		if err != nil {
			c.println("[udp] failed to open UDP socket, trying again in", engineIntervalWatchdog)
			continue
		}

		// start connection watchdog
		connCtx, connCancel := context.WithCancel(ctx)
		go c.connWatchdog(connCtx, conn)
		c.println("[davisweather udp] listening for weather broadcasts")

		// create internal buffer for reading UDP broadcasts
		buff := make([]byte, udpBufferSize)
		for {
			select {
			case <-ctx.Done():
				// terminate connection and UDP event loop
				connCancel()
				return
			default:
			}

			// configure timeouts
			conn.SetReadDeadline(time.Now().Add(udpDeadline))

			// read from UDP
			n, _, err := conn.ReadFrom(buff)
			if err != nil {
				c.println("[davisweather udp] failed to read from UDP socket, reprovisioning")
				// terminate connection, establish new connection
				connCancel()
				break
			}
			// parse UDP broadcast message
			conditions, err := parser.ParseUDP(buff[:n])
			if err != nil {
				c.println("[davisweather udp] failed to parse broadcast")
				continue
			}
			// update Report state
			err = c.report.UpdateUDP(conditions)
			if err != nil {
				c.println("[davisweather udp] failed to update Report", err)
				continue
			}
			c.udpLastReported = time.Now()
		}
	}
}

// connWatchdog listens for the context Done signal and terminates the UDP
// connection.
func (c *Client) connWatchdog(ctx context.Context, conn *net.UDPConn) {
	select {
	case <-ctx.Done():
		conn.Close()
		return
	}
}

// udpWatchdog fetches the broadcast response if no UDP broadcasts are received
// within the UDP deadline. It updates the udpPort in Client when the broadcast
// port is obtained. It calls the resolved cancel function when the UDP port is
// established.
func (c *Client) udpWatchdog(ctx context.Context, resolved context.CancelFunc) {
	eventTimer := time.NewTimer(engineIntervalWatchdog)
	defer eventTimer.Stop()

	c.println("[davisweather udp] initializing watchdog")

	for {
		eventTimer.Reset(engineIntervalWatchdog)

		// calculate age of last UDP broadcast
		delta := time.Now().Sub(c.udpLastReported)
		if delta > udpDeadline {
			// exceeded UDP deadline, must fetch broadcast response
			broadcast, err := c.fetchBroadcastResponse(ctx)
			if err != nil {
				c.println("[davisweather udp] failed to enable UDP broadcasts", err)
			} else {
				// received port, update port in Client and notify resolved port
				previousPort := c.udpPort
				c.udpPort = broadcast.ConnInfo.Port
				if previousPort == 0 {
					resolved()
				}
				delta = time.Duration(broadcast.ConnInfo.Duration) * time.Second
				c.println("[davisweather udp] enabled UDP broadcasts for", delta)
			}
		}
		// terminate or sleep
		select {
		case <-ctx.Done():
			c.println("[davisweather udp] terminating watchdog")
			return
		case <-eventTimer.C:
		}
	}
}

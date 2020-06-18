// Copyright (c) 2020 Tanner Ryan. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

package davisweather

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

var (
	// errInvalidHostname is returned when empty hostname is provided to Client
	errInvalidHostname = errors.New("davisweather: must supply valid IP address or hostname")
)

const (
	// defaultPort is the default port of the WLL unit
	clientDefaultPort = 80
)

// Client is the Davis weather client.
type Client struct {
	Notify <-chan bool // Notify emits a bool when a new weather report is generated

	report          *Report   // report is the weather report state
	verbose         bool      // verbose enables Client logging to stdout
	unit            *wllUnit  // unit contains the network parameters for connecting to WLL unit
	udpPort         int       // udpPort is the port of the UDP broadcasts
	udpLastReported time.Time // udpLastReported is the time the last UDP report was received

	wg *sync.WaitGroup // wg is for checking if all goroutines are done
}

// Managed returns a managed Davis weather client. It accepts a context for
// cancelling the Client. It automatically discovers the WeatherLink Live (WLL)
// unit on the local network. If multicast DNS is disabled on the network, use
// the unmanaged client. It returns a new client for consuming weather data.
// Verbose enables Client logging.
func Managed(ctx context.Context, verbose bool) *Client {
	// initialize report and notification channel
	report, notify := NewReport(verbose)
	// generate client
	c := &Client{
		Notify:  notify,
		report:  report,
		verbose: verbose,
		wg:      &sync.WaitGroup{},
	}
	c.println("[davisweather] managed client initialized")

	// mDNS context to notify engine when to start
	mDNSCtx, mDNSDone := context.WithCancel(ctx)

	// start mDNS discovery and event engine
	c.wg.Add(2)
	go c.discovery(ctx, mDNSDone)
	go c.engine(ctx, mDNSCtx)
	return c
}

// Unmanaged returns an unmanaged Davis weather client. It accepts a context for
// cancelling the Client and a hostname (IP or domain) and port of the
// WeatherLink Live (WLL) unit. It returns a new client for consuming weather
// data. It returns an error if no hostname is provided. Verbose enables Client
// logging.
func Unmanaged(ctx context.Context, verbose bool, hostname string, port int) (*Client, error) {
	if hostname == "" {
		return nil, errInvalidHostname
	}
	// if no port provided, use
	if port <= 0 {
		port = clientDefaultPort
	}
	// initialize report and notification channel
	var u wllUnit
	report, notify := NewReport(verbose)

	// parse provided hostname
	ip := net.ParseIP(hostname)
	switch ip {
	case nil:
		// not an IP address
		u.HostName = hostname
		u.Port = port
	default:
		// IP address
		ipv4 := ip.To4() != nil
		u.HostName = hostname
		u.Port = port
		if ipv4 {
			u.AddrIPv4 = append(u.AddrIPv4, ip)
		} else {
			u.HostName = "[" + u.HostName + "]"
			u.AddrIPv6 = append(u.AddrIPv6, ip)
		}
	}
	// generate client
	c := &Client{
		Notify:  notify,
		report:  report,
		verbose: verbose,
		unit:    &u,
		wg:      &sync.WaitGroup{},
	}
	c.printf("[davisweather] unmanaged client initialized, using WeatherLink Live unit at %s:%d", u.HostName, u.Port)

	// start event engine, no mDNS context
	c.wg.Add(1)
	go c.engine(ctx, nil)
	return c, nil
}

// Report returns the latest weather report or an error.
func (c *Client) Report() (*Report, error) {
	return c.report.Copy()
}

// Closed blocks until the client has been gracefully terminated.
func (c *Client) Closed() {
	c.wg.Wait()
}

// println calls log.Println if verbose logging is enabled.
func (c *Client) println(v ...interface{}) {
	if c.verbose {
		log.Println(v...)
	}
}

// printf calls log.Printf is verbose logging is enabled. Like Println, it
// terminates with a new line.
func (c *Client) printf(format string, v ...interface{}) {
	if c.verbose {
		log.Printf(format+"\n", v...)
	}
}

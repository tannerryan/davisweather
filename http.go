// Copyright (c) 2020 Tanner Ryan. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

package davisweather

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tannerryan/davisweather/parser"
)

var (
	// httpClient is an HTTP client with keep alives disabled
	httpClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
)

// fetchConditionsHTTP fetches weather conditions over HTTP. It returns an error
// if the HTTP response is not formatted correctly, or if the request fails.
func (c *Client) fetchConditionsHTTP(ctx context.Context) (*parser.ConditionsHTTP, error) {
	url := fmt.Sprintf("%s%s", c.unit.GetURL(), routeConditions)
	// prepare request context
	ctx, cancel := context.WithTimeout(ctx, httpTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	// perform request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// parse body
	conditions, err := parser.ParseHTTP(body)
	if err != nil {
		return nil, err
	}
	if conditions.Error != nil {
		return nil, errors.New(conditions.Error.Message)
	}
	return conditions, nil
}

// fetchBroadcastResponse fetches the enable UDP broadcast broadcast response
// over HTTP. It returns an error if the HTTP response is not formatted
// correctly, or if the request fails.
func (c *Client) fetchBroadcastResponse(ctx context.Context) (*parser.BroadcastResponse, error) {
	url := fmt.Sprintf("%s%s?duration=%.0f",
		c.unit.GetURL(), routeBroadcastResponse, udpDuration.Seconds())
	// prepare request context
	ctx, cancel := context.WithTimeout(ctx, httpTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	// perform request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// parse body
	broadcastResp, err := parser.ParseBroadcastResponse(body)
	if err != nil {
		return nil, err
	}
	if broadcastResp.Error != nil {
		return nil, errors.New(broadcastResp.Error.Message)
	}
	return broadcastResp, nil
}

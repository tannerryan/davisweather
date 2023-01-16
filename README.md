# davisweather
[![Build
Status](https://img.shields.io/travis/tannerryan/davisweather.svg?style=flat-square)](https://travis-ci.org/tannerryan/davisweather)
[![Go Report
Card](https://goreportcard.com/badge/github.com/tannerryan/davisweather?style=flat-square)](https://goreportcard.com/report/github.com/tannerryan/davisweather)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5673AF.svg?style=flat-square)](https://pkg.go.dev/github.com/tannerryan/davisweather)
[![GitHub
license](https://img.shields.io/github/license/tannerryan/davisweather.svg?style=flat-square)](https://github.com/tannerryan/davisweather/blob/master/LICENSE)

Real time consumption of weather data from [Davis WeatherLink
Live](https://www.davisinstruments.com/weatherlinklive/). 


## Table of Contents
- [About](#about)
- [Usage](#usage)
    - [Managed Client](#managed-client)
    - [Unmanaged Client](#unmanaged-client)
- [License](#license)


## About
davisweather is used to consume from a WeatherLink Live (WLL) unit located on a
network. It takes advantage of Davis' live transmission protocol, allowing for
updates every 2.5 seconds.


## Usage
A managed and unmanaged client is available. The managed client automatically
discovers the WLL unit on a local network using mDNS. The unmanaged client is
used on networks that block mDNS.

An example client may be found in the [example](example/main.go) directory.

### Managed Client
Here is an example of using the managed client with verbose logging disabled.
```go
func main() {
    ctx := context.Background()
    client := davisweather.Managed(ctx, false)

    for {
        <-client.Notify
        report, err := client.Report()
        if err != nil {
            panic(err)
        }
        log.Println(*report.Temperature)
    }
}
```

### Unmanaged Client
Here is an example of using the unmanaged client with verbose logging enabled.
```go
func main() {
    ctx := context.Background()
    client := davisweather.Unmanaged(ctx, false, "10.0.0.2", 80)

    for {
        <-client.Notify
        report, err := client.Report()
        if err != nil {
            panic(err)
        }
        log.Println(*report.Temperature)
    }
}
```

### Client Shutdown
To shutdown the client, send a Done signal on the context provided to the
client.


## License
Copyright (c) 2020 Tanner Ryan. All rights reserved. Use of this source code is
governed by a BSD-style license that can be found in the LICENSE file.

The ZeroConf package is distributed under an MIT license. Copyright (c) 2016
Stefan Smarzly, Copyright (c) 2014 Oleksandr Lobunets. All rights reserved.

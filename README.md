# trustar-golang
[![Build Status](https://travis-ci.org/jakewarren/trustar-golang.svg?branch=master)](https://travis-ci.org/jakewarren/trustar-golang/)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/jakewarren/trustar-golang)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/jakewarren/trustar-golang/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakewarren/trustar-golang)](https://goreportcard.com/report/github.com/jakewarren/trustar-golang)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=shields)](http://makeapullrequest.com)
> Golang SDK for the TruSTAR API

## API Documentation

See https://docs.trustar.co/api/index.html for the official documentation to the TruSTAR API.

## Example

### Get reports from an enclave
```golang
package main

import (
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	trustar "github.com/jakewarren/trustar-golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Timestamp().Logger()
	
	// set the list of enclave IDs we want to pull from
	enclaveIDs := []string{"abc-123-def"}

	// initialize the client with our API auth information
	c, err := trustar.NewClient("TruSTAR API Key goes here", "TruSTAR API Secret goes here", trustar.APIBaseLive)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating client")
	}

	// acquire the access token
	_, err = c.GetAccessToken()
	if err != nil {
		log.Fatal().Err(err).Msg("error while getting access token")
	}

	t := time.Now().Add(-24 * time.Hour)

	// build the request and send it
	query := url.Values{}
	query.Add("from", strconv.FormatInt(timeToMsEpoch(t), 10)) // request reports from the past 24 hours
	query.Add("enclaveIds", strings.Join(enclaveIDs, ","))
	reports, err := c.GetReports(query)

	// print out the reports we received
	spew.Dump(reports)
}

// convert time.Time to milliseconds epoch string
func timeToMsEpoch(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

```

## Roadmap

Implemented endpoints can be found in [TODO.md](TODO.md)

## Disclaimer

This project is currently in alpha status. Breaking changes are possible and the project has not been thoroughly tested.

## Acknowledgements

The overall design of this client is based upon [logpacker/PayPal-Go-SDK](https://github.com/logpacker/PayPal-Go-SDK)

## License

MIT © 2019 Jake Warren

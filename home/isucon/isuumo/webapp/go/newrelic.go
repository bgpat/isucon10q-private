package main

import (
	"fmt"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var app *newrelic.Application

func init() {
	var err error
	app, err = newrelic.NewApplication(
		newrelic.ConfigAppName("bgpat/isucon10q-private"),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigDebugLogger(os.Stdout),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

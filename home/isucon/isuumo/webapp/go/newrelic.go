package main

import (
	"context"
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v3"
	"github.com/newrelic/go-agent/v3/newrelic"
)

var app *newrelic.Application

func init() {
	var err error
	app, err = newrelic.NewApplication(
		newrelic.ConfigAppName("bgpat/isucon10q-private"),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigDistributedTracerEnabled(true),
		//newrelic.ConfigDebugLogger(os.Stdout),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func nrctx(c echo.Context) context.Context {
	txn := nrecho.FromContext(c)
	ctx := newrelic.NewContext(c.Request().Context(), txn)
	return ctx
}

func nrsgmt(ctx context.Context, name string) *newrelic.Segment {
	txn := newrelic.FromContext(ctx)
	return newrelic.StartSegment(txn, name)
}

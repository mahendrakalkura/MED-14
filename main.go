package main

import (
	"flag"
	"github.com/getsentry/raven-go"
)

func main() {
	action := flag.String("action", "", "")

	flag.Parse()

	settings := get_settings()

	raven.SetDSN(settings.Sentry.Dsn)

	if *action == "bootstrap" {
		bootstrap(settings)
	}

	if *action == "insert" {
		insert(settings)
	}

	if *action == "addresses_one" {
		addresses_one(settings)
	}

	if *action == "addresses_all" {
		addresses_all(settings)
	}

	if *action == "progress" {
		progress(settings)
	}

	if *action == "report" {
		report(settings)
	}
}

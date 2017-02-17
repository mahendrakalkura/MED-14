package main

import (
	"fmt"
	"github.com/getsentry/raven-go"
)

func addresses_one(settings *Settings) {
	fmt.Println("addresses_one()")

	database := get_database(settings)

	address := addresses_select_star_pending_one(database)
	fmt.Printf("%-11s: %s\n", "Street", address.Street)

	results, err := get_results(settings, address)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
	for _, result := range results {
		fmt.Printf("%-11s: %s\n", "Location", result.Location)
		fmt.Printf("%-11s: %s\n", "Office Name", result.OfficeName)
		fmt.Printf("%-11s: %s\n", "Office URL", result.OfficeUrl)
		fmt.Printf("%-11s: %s\n", "Quarter", result.Quarter)
		fmt.Printf("%-11s: %s\n", "Street", result.Street)
	}
}

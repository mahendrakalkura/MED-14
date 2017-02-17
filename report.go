package main

import (
	"encoding/csv"
	"fmt"
	"github.com/getsentry/raven-go"
	"gopkg.in/cheggaaa/pb.v1"
	"os"
)

func report(settings *Settings) {
	fmt.Println("report()")

	database := get_database(settings)

	count := addresses_select_count(database)
	if count == 0 {
		return
	}

	rows := addresses_select_star(database)

	file, create_err := os.Create("Strasse_Zürich_und_Winterthur - Report.csv")
	if create_err != nil {
		raven.CaptureErrorAndWait(create_err, nil)
		panic(create_err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	write_err := writer.Write(
		[]string{
			"Street (input)",
			"Location",
			"Office Name",
			"Office URL",
			"Quarter",
			"Street",
		},
	)
	if write_err != nil {
		raven.CaptureErrorAndWait(write_err, nil)
		panic(write_err)
	}

	progress_bar := pb.StartNew(count)
	for rows.Next() {
		var address_and_result AddressAndResult
		struct_scan_err := rows.StructScan(&address_and_result)
		if struct_scan_err != nil {
			raven.CaptureErrorAndWait(struct_scan_err, nil)
			panic(struct_scan_err)
		}

		write_err := writer.Write(
			[]string{
				address_and_result.AddressStreet,
				get_text(address_and_result.ResultLocation),
				get_text(address_and_result.ResultOfficeName),
				get_text(address_and_result.ResultOfficeUrl),
				get_text(address_and_result.ResultQuarter),
				get_text(address_and_result.ResultStreet),
			},
		)
		if write_err != nil {
			raven.CaptureErrorAndWait(write_err, nil)
			panic(write_err)
		}

		progress_bar.Increment()
	}

	defer writer.Flush()
}

package main

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/jmoiron/sqlx"
	"gopkg.in/cheggaaa/pb.v1"
	"os"
	"os/signal"
	"syscall"
)

func addresses_all(settings *Settings) {
	fmt.Println("addresses_all()")

	signal_channel := make(chan os.Signal)
	addresses_channel := make(chan Address, settings.Others.Consumers*2)

	database := get_database(settings)

	for index := 1; index <= settings.Others.Consumers; index++ {
		go addresses_all_consumer(settings, database, addresses_channel)
	}

	go addresses_all_producer(database, addresses_channel)

	signal.Notify(signal_channel, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	<-signal_channel

	close(addresses_channel)
}

func addresses_all_consumer(settings *Settings, database *sqlx.DB, addresses_channel chan Address) {
	for address := range addresses_channel {
		results, err := get_results(settings, address)
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
		} else {
			address.Status = "Completed"
			addresses_update(database, address, results)
		}
	}
}

func addresses_all_producer(database *sqlx.DB, addresses_channel chan Address) {
	total := addresses_select_count_pending(database)
	rows := addresses_select_star_pending(database)
	progress_bar := pb.StartNew(total)
	for rows.Next() {
		var address Address
		struct_scan_err := rows.StructScan(&address)
		if struct_scan_err != nil {
			raven.CaptureErrorAndWait(struct_scan_err, nil)
		} else {
			addresses_channel <- address
		}
		progress_bar.Increment()
	}
}

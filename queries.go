package main

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/jmoiron/sqlx"
	"strconv"
)

func addresses_select_count(database *sqlx.DB) int {
	statement := `SELECT COUNT(id) FROM addresses`
	row := database.QueryRow(statement)
	var count int
	err := row.Scan(&count)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
	return count
}

func addresses_select_count_pending(database *sqlx.DB) int {
	statement := `SELECT COUNT(id) FROM addresses WHERE status = $1`
	status := "Pending"
	row := database.QueryRow(statement, status)
	var count int
	err := row.Scan(&count)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
	return count
}

func addresses_select_star(database *sqlx.DB) *sqlx.Rows {
	statement := `
	SELECT
		addresses.street AS addresses_street,
		results.location AS results_location,
		results.office_name AS results_office_name,
		results.office_url AS results_office_url,
		results.quarter AS results_quarter,
		results.street AS results_street
	FROM addresses
	LEFT OUTER JOIN results ON results.address_id = addresses.id
	ORDER BY addresses.id ASC, results.id ASC
	`
	rows, queryx_err := database.Queryx(statement)
	if queryx_err != nil {
		raven.CaptureErrorAndWait(queryx_err, nil)
	}
	return rows
}

func addresses_select_star_pending(database *sqlx.DB) *sqlx.Rows {
	statement := `SELECT * FROM addresses WHERE status = $1 ORDER BY id ASC`
	status := "Pending"
	rows, queryx_err := database.Queryx(statement, status)
	if queryx_err != nil {
		raven.CaptureErrorAndWait(queryx_err, nil)
	}
	return rows
}

func addresses_select_star_pending_one(database *sqlx.DB) Address {
	statement := `
    SELECT *
    FROM addresses
    WHERE status = $1
    ORDER BY RANDOM()
    LIMIT 1
    OFFSET 0
    `
	status := "Pending"
	var address Address
	err := database.Get(&address, statement, status)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
	return address
}

func addresses_update(database *sqlx.DB, address Address, results []Result) {
	for _, result := range results {
		statement := `
		INSERT INTO results
		(address_id, location, office_name, office_url, quarter, street)
        VALUES
        (:address_id, :location, :office_name, :office_url, :quarter, :street)
        `
		database.NamedExec(statement, result)
	}

	statement := `UPDATE addresses SET status =: status WHERE id = :id`
	database.NamedExec(statement, address)
}

func addresses_progress(database *sqlx.DB) (string, string, string, string) {
	total := addresses_select_count(database)
	total_string := fmt.Sprintf("%07s", strconv.Itoa(total))

	pending := addresses_select_count_pending(database)
	pending_string := fmt.Sprintf("%07s", strconv.Itoa(pending))

	completed := total - pending
	completed_string := fmt.Sprintf("--%07s", strconv.Itoa(completed))

	percentage := (float64(completed) * 100.00) / (float64(total) * 1.00)
	percentage_string := fmt.Sprintf("---%06.2f%%", percentage)

	return total_string, completed_string, pending_string, percentage_string
}

func results_select_count(database *sqlx.DB, address Address) int {
	statement := `SELECT COUNT(id) FROM results WHERE address_id = $1 ORDER BY id ASC`
	row := database.QueryRow(statement, address.Id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
	return count
}

func results_select_star(database *sqlx.DB, address Address) *sqlx.Rows {
	statement := `SELECT * FROM results WHERE address_id = $1 ORDER BY id ASC`
	rows, queryx_err := database.Queryx(statement, address.Id)
	if queryx_err != nil {
		raven.CaptureErrorAndWait(queryx_err, nil)
	}
	return rows
}

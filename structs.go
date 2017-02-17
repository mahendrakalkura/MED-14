package main

import (
	"database/sql"
)

type Settings struct {
	Proxies SettingsProxies `toml:"proxies"`
	Sentry  SettingsSentry  `toml:"sentry"`
	SQLX    SettingsSQLX    `toml:"sqlx"`
	Others  SettingsOthers  `toml:"others"`
}

type SettingsProxies struct {
	Hostname string `toml:"hostname"`
	Ports    []int  `toml:"ports"`
}

type SettingsSentry struct {
	Dsn string `toml:"dsn"`
}

type SettingsSQLX struct {
	Database string `toml:"database"`
	Hostname string `toml:"hostname"`
	Password string `toml:"password"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
}

type SettingsOthers struct {
	Consumers int `toml:"consumers"`
}

type Address struct {
	Id     int    `db:"id"`
	Street string `db:"street"`
	Status string `db:"status"`
}

type Result struct {
	Id         int    `db:"id"`
	AddressId  int    `db:"address_id"`
	Location   string `db:"location"`
	OfficeName string `db:"office_name"`
	OfficeUrl  string `db:"office_url"`
	Quarter    string `db:"quarter"`
	Street     string `db:"street"`
}

type AddressAndResult struct {
	AddressStreet    string         `db:"addresses_street"`
	ResultLocation   sql.NullString `db:"results_location"`
	ResultOfficeName sql.NullString `db:"results_office_name"`
	ResultOfficeUrl  sql.NullString `db:"results_office_url"`
	ResultQuarter    sql.NullString `db:"results_quarter"`
	ResultStreet     sql.NullString `db:"results_street"`
}

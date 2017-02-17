package main

import (
	"fmt"
)

func bootstrap(settings *Settings) {
	fmt.Println("bootstrap()")

	database := get_database(settings)

	statement := `
    DROP SCHEMA IF EXISTS public CASCADE;

    CREATE SCHEMA IF NOT EXISTS public;

    CREATE TABLE IF NOT EXISTS addresses
    (
        id INTEGER NOT NULL,
        street TEXT NOT NULL,
        status TEXT NOT NULL DEFAULT 'Pending'
    );

    CREATE SEQUENCE addresses_id_sequence;

    ALTER TABLE addresses ALTER COLUMN id SET DEFAULT NEXTVAL
    ('addresses_id_sequence'::regclass);

    ALTER TABLE addresses ADD CONSTRAINT addresses_id_constraint PRIMARY KEY
    (id);

    CREATE INDEX addresses_street ON addresses USING btree (street);

    CREATE INDEX addresses_status ON addresses USING btree (status);

    CREATE TABLE IF NOT EXISTS results
    (
        id INTEGER NOT NULL,
        address_id INTEGER NOT NULL,
        location TEXT NOT NULL,
        office_name TEXT NOT NULL,
        office_url TEXT NOT NULL,
        quarter TEXT NOT NULL,
        street TEXT NOT NULL
    );

    CREATE SEQUENCE results_id_sequence;

    ALTER TABLE results ALTER COLUMN id SET DEFAULT NEXTVAL
    ('results_id_sequence'::regclass);

    ALTER TABLE results ADD CONSTRAINT results_id_constraint PRIMARY KEY (id);

    ALTER TABLE results ADD CONSTRAINT results_address_id
    FOREIGN KEY (address_id) REFERENCES addresses (id) ON DELETE CASCADE
    DEFERRABLE INITIALLY DEFERRED;
    `

	defer database.Close()

	database.MustExec(statement)
}

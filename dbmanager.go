package main

import (
	"context"
	"github.com/jackc/pgx/pgxpool"
)

type Database struct{
	pool *pgxpool.Pool
}

func NewDatabase(connString string) *Database{
	pool, err := pgxpool.Connect(context.Background(), connString)
	checkErr(err)
	d := new(Database)
	d.pool = pool
	return d
}

/**
Generates table if it doesn't yet exist
Useful and convenient when deploying on a new system
*/
func (d *Database) GenerateTable() error {
	execstring := `
CREATE TABLE IF NOT EXISTS userdata (
email text NOT NULL UNIQUE,
sofi_money text NOT NULL DEFAULT '',
sofi_money_clicks integer NOT NULL DEFAULT 0,
sofi_invest text NOT NULL DEFAULT '',
sofi_invest_clicks integer NOT NULL DEFAULT 0,
robinhood text NOT NULL DEFAULT '',
robinhood_clicks integer NOT NULL DEFAULT 0,
amazon text NOT NULL DEFAULT '',
amazon_clicks integer NOT NULL DEFAULT 0,
airbnb text NOT NULL DEFAULT '',
airbnb_clicks integer NOT NULL DEFAULT 0,
grubhub text NOT NULL DEFAULT '',
grubhub_clicks integer NOT NULL DEFAULT 0,
doordash text NOT NULL DEFAULT '',
doordash_clicks integer NOT NULL DEFAULT 0,
uber text NOT NULL DEFAULT '',
uber_clicks integer NOT NULL DEFAULT 0
);
`
	_, err := d.pool.Exec(context.Background(), execstring)
	return err
}
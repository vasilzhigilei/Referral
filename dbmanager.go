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

func (d *Database) GetServiceURLs(service string) []string {
	querystring := "SELECT " + service + " FROM userdata WHERE " + service + " != '';"
	rows, err := d.pool.Query(context.Background(), querystring)
	checkErr(err)
	var returnvalue []string
	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		returnvalue = append(returnvalue, temp)
	}
	return returnvalue
}

func (d *Database) GetUser(email string) *User {
	querystring := "SELECT * FROM userdata WHERE email = '" + email + "';"
	rows, err := d.pool.Query(context.Background(), querystring)
	checkErr(err)
	user := User{}
	for rows.Next() {
		err = rows.Scan(&user.Email, &user.Sofi_money, &user.Sofi_money_clicks, &user.Sofi_invest, &user.Sofi_invest_clicks,
			&user.Robihood, &user.Robinhood_clicks, &user.Amazon, &user.Amazon_clicks, &user.Airbnb, &user.Airbnb_clicks, &user.Grubhub,
			&user.Grubhub_clicks, &user.Doordash, &user.Doordash_clicks, &user.Uber, &user.Uber_clicks)
		checkErr(err)
	}
	return &user
}

/**
Insert user, if conflict (user already exists), do nothing
*/
func (d *Database) InsertUser(email string) error {
	_, err := d.pool.Exec(context.Background(), "INSERT INTO userdata values($1) ON CONFLICT DO NOTHING", email)
	return err
}
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
Email text NOT NULL UNIQUE,
Sofi_money text NOT NULL DEFAULT '',
Sofi_money_clicks integer NOT NULL DEFAULT 0,
Sofi_invest text NOT NULL DEFAULT '',
Sofi_invest_clicks integer NOT NULL DEFAULT 0,
Robinhood text NOT NULL DEFAULT '',
Robinhood_clicks integer NOT NULL DEFAULT 0,
Amazon text NOT NULL DEFAULT '',
Amazon_clicks integer NOT NULL DEFAULT 0,
Airbnb text NOT NULL DEFAULT '',
Airbnb_clicks integer NOT NULL DEFAULT 0,
Grubhub text NOT NULL DEFAULT '',
Grubhub_clicks integer NOT NULL DEFAULT 0,
Doordash text NOT NULL DEFAULT '',
Doordash_clicks integer NOT NULL DEFAULT 0,
Uber text NOT NULL DEFAULT '',
Uber_clicks integer NOT NULL DEFAULT 0
);
`
	_, err := d.pool.Exec(context.Background(), execstring)
	return err
}

func (d *Database) GetServiceURLs(service string) []EmailURLPair {
	querystring := "SELECT email, " + service + " FROM userdata WHERE " + service + " != '';"
	rows, err := d.pool.Query(context.Background(), querystring)
	checkErr(err)
	var returnvalue []EmailURLPair
	for rows.Next() {
		var temp EmailURLPair
		err = rows.Scan(&temp.Email, &temp.URL)
		returnvalue = append(returnvalue, temp)
	}
	return returnvalue
}

func (d *Database) GetUser(email string) *User {
	querystring := "SELECT * FROM userdata WHERE Email = '" + email + "';"
	rows, err := d.pool.Query(context.Background(), querystring)
	checkErr(err)
	user := User{}
	for rows.Next() {
		err = rows.Scan(&user.Email, &user.Sofi_money, &user.Sofi_money_clicks, &user.Sofi_invest, &user.Sofi_invest_clicks,
			&user.Robinhood, &user.Robinhood_clicks, &user.Amazon, &user.Amazon_clicks, &user.Airbnb, &user.Airbnb_clicks,
			&user.Grubhub, &user.Grubhub_clicks, &user.Doordash, &user.Doordash_clicks, &user.Uber, &user.Uber_clicks)
		checkErr(err)
	}
	return &user
}

func (d *Database) UpdateUser(user *User) error {
	execstring := `
UPDATE userdata
SET Sofi_money = $1, Sofi_invest = $2, Robinhood = $3, Amazon = $4, Airbnb = $5, Grubhub = $6, Doordash = $7, Uber = $8 
WHERE email = $9;
`
	_, err := d.pool.Exec(context.Background(), execstring, user.Sofi_money, user.Sofi_invest, user.Robinhood,
		user.Amazon, user.Airbnb, user.Grubhub, user.Doordash, user.Uber, user.Email)

	return err
}

/**
Insert user, if conflict (user already exists), do nothing
*/
func (d *Database) InsertUser(email string) error {
	_, err := d.pool.Exec(context.Background(), "INSERT INTO userdata values($1) ON CONFLICT DO NOTHING", email)
	return err
}
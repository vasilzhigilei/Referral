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
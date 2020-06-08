package main

import (
	"context"
	"database/sql"
	"flag"
	"time"

	_ "github.com/lib/pq"
	"potentialtech.com/txmanager"
)

func main() {
	cs := flag.String("c", "", "Postgres connection string")
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := sql.Open("postgres", *cs)
	if err != nil {
		panic(err)
	}
	pgf := txmanager.MakePostgresTxFinalizer(ctx, "main", pool)
	_, err = pgf.TX.ExecContext(ctx, "CREATE TABLE temp (id INT)")
	if err != nil {
		panic(err)
	}
	_, err = pgf.TX.ExecContext(ctx, "INSERT INTO TEMP (id) VALUES (5)")
	if err != nil {
		panic(err)
	}
	err = pgf.Finalize()
	if err != nil {
		panic(err)
	}
	pgf.Commit()
}

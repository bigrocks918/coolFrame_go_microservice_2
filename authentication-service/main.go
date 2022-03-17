package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
	"time"
)

var conn *sql.DB
var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

	connectToDB()

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    ":80",
		Handler: app.routes(),
	}
	fmt.Println("Starting authentication end service on port 80")
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func connectToDB() {
	// connect to postgres
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			fmt.Println("Postgres not ready...")
			counts++
		} else {
			conn = connection
			break
		}

		if counts > 10 {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

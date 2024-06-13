package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP сетевой адрес")
	dsn := flag.String("dsn", "snippetbox.db", "Путь к SQLite3 базе данных")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// open sqlite 3 db
	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag
	db, err := sql.Open("sqlite3", *dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// create app for loggin
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// create server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// start server
	infoLog.Printf("Start server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection poo
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

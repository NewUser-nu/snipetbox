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

	// Открытие базы данных SQLite3
	db, err := sql.Open("sqlite3", *dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Создание экземпляра application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Создание сервера HTTP
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(), // Использование маршрутов из приложения
	}

	// Запуск сервера
	infoLog.Printf("Start server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

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

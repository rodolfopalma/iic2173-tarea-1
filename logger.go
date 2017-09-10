package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFileName       = "./requests.db"
	createTableQuery = `CREATE TABLE IF NOT EXISTS requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		method TEXT NOT NULL,
		url TEXT NOT NULL,
		remoteAddress TEXT NOT NULL,
		datetime DATETIME NOT NULL
	);`
	insertRecordQuery = `INSERT INTO
		requests(method, url, remoteAddress, datetime) values(?, ?, ?, datetime('now'))
	;`
	retrieveLast10Requests = `SELECT method, url, remoteAddress, datetime
	from requests ORDER BY datetime DESC LIMIT 10;`
)

type requestRecord struct {
	Method     string
	URL        string
	RemoteAddr string
	Datetime   time.Time
}

func rootHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize the request record.
		newRecord := requestRecord{
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Now(), // NB: no se usa este valor.
		}
		// Create a new database record.
		statement, err := db.Prepare(insertRecordQuery)
		if err != nil {
			panic(err)
		}
		statement.Exec(newRecord.Method, newRecord.URL, newRecord.RemoteAddr)
		// Get the last 10 records.
		rows, err := db.Query(retrieveLast10Requests)
		if err != nil {
			panic(err)
		}
		// NB: se puede optimizar usando array de largo 10 (?).
		var records []*requestRecord
		for rows.Next() {
			record := new(requestRecord)
			rows.Scan(&record.Method, &record.URL, &record.RemoteAddr, &record.Datetime)
			fmt.Println(record.Datetime)
			records = append(records, record)
		}
		// Return HTML template.
		tmpl := template.New("root.html")
		tmpl, err = tmpl.ParseFiles("root.html")
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(w, records)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	// Initialize the database.
	db, err := sql.Open("sqlite3", dbFileName)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	statement, err := db.Prepare(createTableQuery)
	if err != nil {
		panic(err)
	}
	statement.Exec()
	// Set up the handler function, passing a pointer to the database connection.
	http.HandleFunc("/", rootHandler(db))
	// Listen and serve.
	fmt.Println("Listening and serving from :8080...")
	http.ListenAndServe(":8080", nil)
}

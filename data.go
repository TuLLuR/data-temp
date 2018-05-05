package main

import (
	"net/http"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"io"
)

var db *sql.DB
var err error

func index(w http.ResponseWriter, r *http.Request)  {
	_, err := io.WriteString(w, "at index")
	check(err)
}

func users(w http.ResponseWriter, r *http.Request)  {
	rows, err := db.Query(`SELECT uName FROM test;`)
	check(err)
	defer rows.Close()

	var s, name string

	s = "Retrieved records:\n"
	for rows.Next() {
		err = rows.Scan(&name)
		check(err)
		s += name + "\n"
	}
	fmt.Fprintln(w, s)
}

func create(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`CREATE TABLE customer (name VARCHAR(20));`)
	check(err)
	defer stmt.Close()

	rows, err := stmt.Exec()
	check(err)

	n, err := rows.RowsAffected()
	check(err)

	fmt.Fprintln(w, "Created Table 'customer'", n)
}

func insert(w http.ResponseWriter, r *http.Request)  {
	stmt, err := db.Prepare(`INSERT INTO customer VALUES ("James");`)
	check(err)
	defer stmt.Close()

	rows, err := stmt.Exec()
	check(err)

	n, err := rows.RowsAffected()
	check(err)

	fmt.Fprintln(w, "Inserted Record", n)
}

func read(w http.ResponseWriter, r *http.Request)  {
	rows, err := db.Query(`SELECT * FROM customer ;`)
	check(err)
	defer rows.Close()

	var name string
	for rows.Next() {
		err = rows.Scan(&name)
		check(err)
		fmt.Fprintln(w, "Retrieved record:", name)
	}

}

func update(w http.ResponseWriter, r *http.Request)  {
	stmt, err := db.Prepare(`UPDATE customer SET name="Jimmy" WHERE name="James";`)
	check(err)
	defer stmt.Close()

	rows, err := stmt.Exec()
	check(err)

	n, err := rows.RowsAffected()
	check(err)

	fmt.Fprintln(w, "Updated record", n)
}

func del(w http.ResponseWriter, r *http.Request)  {
	stmt, err := db.Prepare(`DELETE FROM customer WHERE name="Jimmy";`)
	check(err)
	defer stmt.Close()

	rows, err := stmt.Exec()
	check(err)

	n, err := rows.RowsAffected()
	check(err)

	fmt.Fprintln(w, "Deleted Record", n)
}

func drop(w http.ResponseWriter, r *http.Request)  {
	stmt, err := db.Prepare(`DROP TABLE customer;`)
	check(err)
	defer stmt.Close()

	_, err = stmt.Exec()
	check(err)

	fmt.Fprintln(w, "Droped table customer")
}

func check(err error)  {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// localhost:3030/dbname?useSsl=flase
	// root@localhost:root@tcp(localhost:3306)/test
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3307)/test")
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	http.HandleFunc("/", index)
	http.HandleFunc("/users", users)
	http.HandleFunc("/create", create)
	http.HandleFunc("/insert", insert)
	http.HandleFunc("/read", read)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", del)
	http.HandleFunc("/drop", drop)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":1012", nil))
}

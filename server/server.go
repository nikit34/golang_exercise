package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

type Change struct {
	Timestamp uint64 `json:"timestamp"`
	Value float32 `json:"value"`
}

// type Changes struct {
// 	Total uint32 `json:"total"`
// 	History []Change `json:"history"`
// }




func view_last_change(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, err := db.Query("SELECT timestamp, value from changes")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var change Change
	var changes []Change
	for result.Next() {
		err := result.Scan(&change.Timestamp, &change.Value)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(change.Timestamp)
		changes = append(changes, change)
	}
	json.NewEncoder(w).Encode(changes)
}

func main() {
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/demo")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/api/btcusdt", view_last_change)
	http.ListenAndServe(":8080", nil)
}
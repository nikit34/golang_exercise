package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

type Change struct {
	Timestamp uint64 `json:"timestamp"`
	Value float32 `json:"value"`
}

type Changes struct {
	Total uint32 `json:"total"`
	History []Change `json:"history"`
}




func view_btcusdt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var changes Changes

	result_total, err := db.Query("SELECT COUNT(*) FROM changes")
	if err != nil {
		panic(err.Error())
	}
	for result_total.Next() {
		err := result_total.Scan(&changes.Total)
		if err != nil {
			panic(err.Error())
		}
	}
	defer result_total.Close()

	switch r.Method {
	case "GET":
		result_history, err := db.Query("SELECT timestamp, value from changes ORDER BY timestamp DESC LIMIT 1")
		if err != nil {
			panic(err.Error())
		}
		defer result_history.Close()

		var change Change
		for result_history.Next() {
			err := result_history.Scan(&change.Timestamp, &change.Value)
			if err != nil {
				panic(err.Error())
			}
			changes.History = append(changes.History, change)
		}
	case "POST":
		r.ParseForm()
		start_time := r.Form.Get("start_time")
    	end_time := r.Form.Get("end_time")
		start_i, start_err := strconv.Atoi(start_time)
		end_i, last_err := strconv.Atoi(end_time)
		if start_err != nil || last_err != nil {
			fmt.Println("Failed to convert filters by time")
		}

		if end_time != "" && start_time != "" && start_i < end_i {
			result_history, err := db.Query(
				"SELECT timestamp, value FROM changes WHERE timestamp > ? AND timestamp < ?", start_time, end_time,
			)
			if err != nil {
				panic(err.Error())
			}
			defer result_history.Close()

			var change Change
			for result_history.Next() {
				err := result_history.Scan(&change.Timestamp, &change.Value)
				if err != nil {
					panic(err.Error())
				}
				changes.History = append(changes.History, change)
			}
		} else if start_i >= end_i {
			fmt.Println("`start_time` must be less than `end_time`")
		} else {
			fmt.Println("For PORT request variables `start_time` and `end_time` should be used")
		}
	}

	json.NewEncoder(w).Encode(changes)
}

func main() {
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/demo")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/api/btcusdt", view_btcusdt)
	http.ListenAndServe(":8080", nil)
}
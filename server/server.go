package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"math"

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

func make_post_response(changes *Changes, result_history *sql.Rows){
	var change Change
	for result_history.Next() {
		err := result_history.Scan(&change.Timestamp, &change.Value)
		if err != nil {
			panic(err.Error())
		}
		changes.History = append(changes.History, change)
	}
}

func make_post_select(changes *Changes) {
	result_history, err := db.Query("SELECT timestamp, value FROM changes")
	if err != nil {
		panic(err.Error())
	}
	defer result_history.Close()

	make_post_response(changes, result_history)
}

func make_post_select_time(changes *Changes, start_time, end_time string) {
	result_history, err := db.Query(
		"SELECT timestamp, value FROM changes WHERE timestamp > ? AND timestamp < ?", start_time, end_time,
	)
	if err != nil {
		panic(err.Error())
	}
	defer result_history.Close()

	make_post_response(changes, result_history)
}

func make_post_select_pagination(changes *Changes, limit, offset string) {
	result_history, err := db.Query(
		"SELECT timestamp, value FROM changes LIMIT ? OFFSET ?", limit, offset,
	)
	if err != nil {
		panic(err.Error())
	}
	defer result_history.Close()

	make_post_response(changes, result_history)
}

func make_post_select_time_pagination(changes *Changes, start_time, end_time, limit, offset string) {
	result_history, err := db.Query(
		"SELECT timestamp, value FROM changes WHERE timestamp > ? AND timestamp < ? LIMIT ? OFFSET ?", start_time, end_time, limit, offset,
	)
	if err != nil {
		panic(err.Error())
	}
	defer result_history.Close()

	make_post_response(changes, result_history)
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
		pagination := r.Form.Get("pagination")

		start_i, start_err := strconv.Atoi(start_time)
		end_i, end_err := strconv.Atoi(end_time)
		pagination_i, pagination_err := strconv.Atoi(pagination)
		num_pages := math.Ceil(float64(changes.Total) / float64(pagination_i))

		if start_err != nil && end_err != nil && pagination_err != nil {
			make_post_select(&changes)
		} else if start_err != nil && strings.Contains(start_err.Error(), "") && end_err != nil && strings.Contains(end_err.Error(), "") && pagination_err == nil {
			for i := 0; i < int(num_pages); i++ {
				make_post_select_pagination(&changes, pagination, string(i * pagination_i))
			}
		} else if start_err == nil && end_err == nil && start_i < end_i && pagination_err != nil && strings.Contains(pagination_err.Error(), "") {
			make_post_select_time(&changes, start_time, end_time)
		} else if start_err == nil && end_err == nil && pagination_err == nil {
			if start_i < end_i {
				for i := 0; i < int(num_pages); i++ {
					make_post_select_time_pagination(&changes, start_time, end_time, pagination, string(i * pagination_i))
				}
			} else {
				panic(fmt.Errorf("`start_time` must be less than `end_time`"))
			}
		} else {
			panic(fmt.Errorf("failed to convert filters"))
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
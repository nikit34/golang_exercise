package main

import (
	"encoding/json"
	"net/http"
)


type HistoryItem struct {
	Timestamp uint64 `json:"timestamp"`
	Value float32 `json:"value"`
}

type Changes struct {
	Total uint32 `json:"total"`
	History []HistoryItem `json:"history_item"`
}


func view_last_change(w http.ResponseWriter, r *http.Request) {
	change := &Changes{
		Total: 1232,
		History: []HistoryItem{
			{
				Timestamp: 12321434,
				Value: 123.45,
			},
			{
				Timestamp: 12321434,
				Value: 123.45,
			},
		},
	}
	json.NewEncoder(w).Encode(change)
}

func main() {
	http.HandleFunc("/api/btcusdt", view_last_change)
	http.ListenAndServe(":8080", nil)
}
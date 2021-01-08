package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aren55555/coindesk-async/coindesk"
)

var (
	c *coindesk.Checker
)

type jsonError struct {
	Error string `json:"error"`
}

func main() {
	c = coindesk.New()

	http.HandleFunc("/start", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Start()
	}))

	http.HandleFunc("/stop", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Stop()
	}))

	http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var serializable interface{}
		serializable, err := c.GetValue()
		if err != nil {
			serializable = jsonError{err.Error()}
		}
		data, _ := json.MarshalIndent(serializable, "", "  ")
		w.Write(data)
	}))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

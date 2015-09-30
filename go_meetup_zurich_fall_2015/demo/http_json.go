package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Bar struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func myJsonHandler(w http.ResponseWriter, r *http.Request) {
	b := &Bar{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(b)
	if err != nil {
		log.Fatal(err)
	}

	res, _ := json.Marshal(b)

	w.Write(res)
}

func main() {
	http.HandleFunc("/foo", myJsonHandler)
	http.ListenAndServe("localhost:8080", nil)
}

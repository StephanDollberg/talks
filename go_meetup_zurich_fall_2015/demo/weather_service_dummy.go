package main

import (
	"math/rand"
	"net/http"
	"time"
)

func temp(w http.ResponseWriter, r *http.Request) {
	timeout := rand.Intn(600)
	time.Sleep(time.Millisecond * time.Duration(timeout))
	w.Write([]byte(`{"temp": 23}`))
}

func forecast(w http.ResponseWriter, r *http.Request) {
	timeout := rand.Intn(600)
	time.Sleep(time.Millisecond * time.Duration(timeout))
	w.Write([]byte(`{"temp": 17}`))
}

func main() {
	http.HandleFunc("/temp", temp)
	http.HandleFunc("/forecast", forecast)
	http.ListenAndServe("localhost:20000", nil)

}

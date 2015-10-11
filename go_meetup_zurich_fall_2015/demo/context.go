package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}

type Result struct {
	Value int
	Err   error
}

func GetForecast(ctx context.Context, city string, out chan Result) {
	req, err := http.NewRequest("GET", "http://localhost:20000/forecast", nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Set("q", city)
	req.URL.RawQuery = q.Encode()
	result := Result{}

	result.Err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var data struct {
			Temp int `json:"temp"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		result.Value = data.Temp
		return nil
	})

	out <- result
	fmt.Println("gr done", result.Err)
}
func GetTemp(ctx context.Context, city string, out chan Result) {
	req, err := http.NewRequest("GET", "http://localhost:20000/temp", nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Set("q", city)
	req.URL.RawQuery = q.Encode()
	result := Result{}

	result.Err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var data struct {
			Temp int `json:"temp"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		result.Value = data.Temp
		return nil
	})

	out <- result
	fmt.Println("gr done", result.Err)
}

func AskWeatherService(ctx context.Context) (int, int, error) {
	city, _ := CityFromContext(ctx)
	tempChan := make(chan Result, 1)
	forecastChan := make(chan Result, 1)

	go GetTemp(ctx, city, tempChan)
	go GetForecast(ctx, city, forecastChan)
	temp := Result{}
	forecast := Result{}

	for i := 0; i < 2; i++ {
		select {
		case temp = <-tempChan:
			fmt.Println("temp ")
			if temp.Err != nil {
				fmt.Println("temp failed")
				return 0, 0, temp.Err
			}
		case forecast = <-forecastChan:
			fmt.Println("forecast ")
			if forecast.Err != nil {
				fmt.Println("forecast failed")
				return 0, 0, forecast.Err
			}
		case <-ctx.Done():
			fmt.Println("cancel")
			return 0, 0, ctx.Err()
		}
	}

	return temp.Value, forecast.Value, nil
}

func NewCityContext(parent context.Context, city string) context.Context {
	return context.WithValue(parent, "city", city)
}

func CityFromContext(parent context.Context) (string, bool) {
	city, ok := parent.Value("city").(string)
	return city, ok
}

func ParseCity(r *http.Request) (string, error) {
	res := r.FormValue("city")

	if res == "" {
		return "", errors.New("invalid city")
	}

	return res, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer fmt.Println("done")
	city, err := ParseCity(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*3000)
	defer cancel()

	ctx = NewCityContext(ctx, city)

	temp, forecast, err := AskWeatherService(ctx)

	if err != nil {
		http.Error(w, "Error retrieving weather", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Current temp is: %d\nForecast is: %d", temp, forecast)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8080", nil)
}

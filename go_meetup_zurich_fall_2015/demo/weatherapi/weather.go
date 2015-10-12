package weather

import (
	"encoding/json"
	"log"
	"net/http"

	"go_meetup_zurich_fall_2015/demo/cityapi"

	"golang.org/x/net/context"
)

const serviceBaseUrl = "http://localhost:20000/"

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

type resultWrapper struct {
	Value int
	Err   error
}

func handleResponse(resp *http.Response, err error, res *resultWrapper) error {
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
	res.Value = data.Temp
	return nil
}

func getImpl(ctx context.Context, reqType string, city string, out chan resultWrapper) {
	req, err := http.NewRequest("GET", serviceBaseUrl+reqType, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Set("q", city)
	req.URL.RawQuery = q.Encode()
	res := resultWrapper{}

	res.Err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		return handleResponse(resp, err, &res)
	})

	out <- res
}

func getForecast(ctx context.Context, city string, out chan resultWrapper) {
	getImpl(ctx, "forecast", city, out)
}

func getTemp(ctx context.Context, city string, out chan resultWrapper) {
	getImpl(ctx, "temp", city, out)
}

type QueryResult struct {
	Temperature int
	Forecast    int
}

func Query(ctx context.Context) (*QueryResult, error) {
	city, _ := city.FromContext(ctx) // HL
	tempChan, forecastChan := make(chan resultWrapper, 1), make(chan resultWrapper, 1)
	go getTemp(ctx, city, tempChan)
	go getForecast(ctx, city, forecastChan)
	temperature, forecast := resultWrapper{}, resultWrapper{}

	for i := 0; i < 2; i++ {
		select {
		case temperature = <-temperatureChan: // HL
			if temperature.Err != nil { // HL
				return nil, temperature.Err
			}
		case forecast = <-forecastChan:
			if forecast.Err != nil {
				return nil, forecast.Err
			}
		case <-ctx.Done(): // HL
			return nil, ctx.Err() // HL
		}
	}
	return &QueryResult{temperature.Value, forecast.Value}, nil
}

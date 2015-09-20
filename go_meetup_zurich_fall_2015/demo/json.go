package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {

	{
		raw_data := []byte(`{"msg": "Hello Go!"}`)
		decoded := map[string]string{}
		err := json.Unmarshal(raw_data, &decoded)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(decoded["msg"])
	}

	{
		raw_data := []byte(`{"msg": "Hello Go!", "id": 12345}`)
		decoded := map[string]interface{}{}
		err := json.Unmarshal(raw_data, &decoded)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(decoded["msg"].(string), int(decoded["id"].(float64)))
	}
}

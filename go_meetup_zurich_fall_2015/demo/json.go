package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Msg1 struct {
	msg string
	id  int64
}

type Msg2 struct {
	Msg string
	Id  int64
}

type Msg3 struct {
	Msg string `json:"msg"`
	Id  int64  `json:"id"`
}

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

	{
		raw_data := []byte(`{"msg": "Hello Go!", "id": 12345}`)
		decoded := Msg1{}
		err := json.Unmarshal(raw_data, &decoded)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(decoded.msg, decoded.id)
	}

	{
		raw_data := []byte(`{"msg": "Hello Go!", "id": 12345}`)
		decoded := Msg2{}
		err := json.Unmarshal(raw_data, &decoded)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(decoded.Msg, decoded.Id)

		encoded, _ := json.Marshal(&decoded)
		fmt.Println(string(encoded))
	}

	{
		raw_data := []byte(`{"msg": "Hello Go!", "id": 12345}`)
		decoded := Msg3{}
		err := json.Unmarshal(raw_data, &decoded)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(decoded.Msg, decoded.Id)
	}
}

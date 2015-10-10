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
	Msg string // HL
	Id  int64  // HL
}

type Msg3 struct {
	Msg string `json:"msg"` // HL
	Id  int64  `json:"id"`  // HL
}

type OmitEmptyStruct struct {
	Msg string `json:"msg,omitempty"` // HL
	Id  int64  `json:"id"`
}

type RawMessageStruct struct {
	Type     int64           `json:"type"`
	UserData json.RawMessage `json:"userdata"` // HL
}

func StringStringMap() {
	raw_data := []byte(`{"msg": "Hello Go!"}`)
	decoded := map[string]string{}
	err := json.Unmarshal(raw_data, &decoded)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(decoded["msg"])
}

func StringInterfaceMap() {
	raw_data := []byte(`{"msg": "Hello Go!", "id": 12345}`)
	decoded := map[string]interface{}{} // HL
	err := json.Unmarshal(raw_data, &decoded)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(decoded["msg"].(string), int(decoded["id"].(float64))) // HL
}

func UnexportedStruct() {
	raw_data := []byte(`{"msg": "Hello Go!", "id": 12345}`)
	decoded := Msg1{}
	err := json.Unmarshal(raw_data, &decoded)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(decoded.msg, decoded.id)
}

func ExportedStruct() {
	raw_data := []byte(`{"msg": "Hello Go!", "id": 12345}`)
	decoded := Msg2{}
	err := json.Unmarshal(raw_data, &decoded)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(decoded.Msg, decoded.Id)

	encoded, _ := json.Marshal(&decoded) // HL
	fmt.Println(string(encoded))
}

func WithTags() {
	raw_data := []byte(`{"msg": "Hello Go!", "id": 12345}`)
	decoded := Msg3{}
	err := json.Unmarshal(raw_data, &decoded)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(decoded.Msg, decoded.Id)
}

func AnonymousStruct() {
	raw_data := []byte(`{"msg": "Hello Go!", "id": 12345}`)
	decoded := struct { // HL
		Msg string `json:"msg"` // HL
		Id  int64  `json:"id"`  // HL
	}{} // HL
	err := json.Unmarshal(raw_data, &decoded)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(decoded.Msg, decoded.Id)
}

func OmitEmpty() {
	data := OmitEmptyStruct{Id: 123}

	encoded, err := json.Marshal(&data)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(encoded))
}

func RawMessage() {
	raw_data := []byte(`{"type": 1, "userdata": {"foo": 1, "bar": 2}}`)

	decoded := RawMessageStruct{}
	err := json.Unmarshal(raw_data, &decoded)

	if err != nil {
		log.Fatal(err)
	}

	// depending on type ...

	userdata := make(map[string]int)
	_ = json.Unmarshal(decoded.UserData, &userdata) // HL
	fmt.Println(userdata)
}

func main() {

	fmt.Println("map[string]string{}")
	StringStringMap()

	fmt.Println("map[string]interface{}")
	StringInterfaceMap()

	fmt.Println("unexported struct")
	UnexportedStruct()

	fmt.Println("exported struct")
	ExportedStruct()

	fmt.Println("exported struct with tags")
	WithTags()

	fmt.Println("annonymous struct with tags")
	AnonymousStruct()

	fmt.Println("omit empty")
	OmitEmpty()

	fmt.Println("RawMessage")
	RawMessage()
}

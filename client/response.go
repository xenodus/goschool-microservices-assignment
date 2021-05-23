package main

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (res *Response) print() {
	fmt.Println("========================================")
	fmt.Println("Status:", res.Status)
	fmt.Println("Code:", res.Code)
	fmt.Println("Message:", res.Message)
	fmt.Println("========================================")
}

func printResponse(data []byte, header string) {
	var res Response
	marshalErr := json.Unmarshal(data, &res)

	if marshalErr != nil {
		fmt.Println("Error decoding json")
	}

	if header != "" {
		fmt.Println("========================================")
		fmt.Println(header)
	}

	res.print()
}

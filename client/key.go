package main

import "fmt"

var apiKey string = ""

func logout() {
	apiKey = ""
}

func printKey() {
	fmt.Println("API Key:", apiKey)
}

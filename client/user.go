package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type User struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type RegistrationResponse struct {
	Status       string `json:"status"`
	Code         int    `json:"code"`
	Apikey       string `json:"apiKey"`
	Apikeystatus string `json:"apiKeyStatus"`
}

func (res *RegistrationResponse) print() {
	fmt.Println("========================================")
	fmt.Println("Status:", res.Status)
	fmt.Println("Code:", res.Code)
	fmt.Println("Apikey:", res.Apikey)
	fmt.Println("========================================")
}

// register and get api key - ideally would include some verification method e.g. user confirming via link to email
func (user *User) register() {

	jsonValue, _ := json.Marshal(user)
	url := apiBaseUrl + "users?apiKey=" + apiKey
	//fmt.Println(url)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode == http.StatusCreated {
			var res RegistrationResponse
			marshalErr := json.Unmarshal(data, &res)

			if marshalErr != nil {
				fmt.Println("Error decoding json")
			}

			res.print()
		} else {
			printResponse(data, "> Error encountered")
		}
	}
}

// user is a user flagged as admin
func (user *User) invalidateKey(key string) {

	jsonValue, _ := json.Marshal(user)
	url := apiBaseUrl + "keys/" + key + "?apiKey=" + apiKey
	//fmt.Println(url)
	request, _ := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode == http.StatusOK {
			printResponse(data, "Success")
		} else {
			printResponse(data, "> Error encountered")
		}
	}
}

func validateEmail(email string) error {

	if len(email) < emailMinLen || len(email) > emailMaxLen {
		return errAuthEmailLength
	}

	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !emailRegex.MatchString(email) {
		return errAuthEmailFormat
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < passwordMinLen || len(password) > passwordMaxLen {
		return errAuthPasswordLength
	}

	return nil
}

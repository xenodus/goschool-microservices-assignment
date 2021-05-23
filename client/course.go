package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Course struct {
	Courseid    string `json:"CourseId"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Status      string `json:"Status"`
}

func getCourses(courseType string) {

	var url string

	switch courseType {
	case "all":
		url = apiBaseUrl + "courses/all?apiKey=" + apiKey
	case "inactive":
		url = apiBaseUrl + "courses/inactive?apiKey=" + apiKey
	default:
		url = apiBaseUrl + "courses?apiKey=" + apiKey
	}

	//fmt.Println(url)
	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		if response.StatusCode == http.StatusOK {
			courses := []Course{}
			marshalErr := json.Unmarshal(data, &courses)

			if marshalErr != nil {
				fmt.Println("Error decoding json")
			}

			if len(courses) > 0 {

				fmt.Println("========================================")

				if courseType == "all" {
					fmt.Println("> Showing all", len(courses), "courses:")
				} else {
					fmt.Println("> Showing all", len(courses), courseType, "courses:")
				}

				for _, course := range courses {
					fmt.Println("========================================")
					fmt.Println("ID:", course.Courseid)
					fmt.Println("Title:", course.Title)
					fmt.Println("Description:", course.Description)
					fmt.Println("Status:", course.Status)
					fmt.Println("========================================")
				}
			} else {
				fmt.Println("========================================")
				fmt.Println("> No courses found")
				fmt.Println("========================================")
			}
		} else {
			printResponse(data, "> Error encountered")
		}
	}
}

func getCourse(courseId string) {
	url := apiBaseUrl + "courses/" + courseId + "?apiKey=" + apiKey
	//fmt.Println(url)
	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {

		data, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		if response.StatusCode == http.StatusOK {
			var course Course
			marshalErr := json.Unmarshal(data, &course)

			if marshalErr != nil {
				fmt.Println("Error decoding json")
			}

			fmt.Println("========================================")
			fmt.Println("> Showing information for course:", course.Courseid)
			fmt.Println("========================================")
			fmt.Println("ID:", course.Courseid)
			fmt.Println("Title:", course.Title)
			fmt.Println("Description:", course.Description)
			fmt.Println("Status:", course.Status)
			fmt.Println("========================================")
		} else {
			printResponse(data, "> Error encountered")
		}
	}
}

func (course *Course) createCourse() {

	jsonValue, _ := json.Marshal(course)
	url := apiBaseUrl + "courses/" + course.Courseid + "?apiKey=" + apiKey
	//fmt.Println(url)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode == http.StatusCreated {
			printResponse(data, "Success")
		} else {
			printResponse(data, "> Error encountered")
		}
	}
}

func (course *Course) updateCourse(courseId string) {

	jsonValue, _ := json.Marshal(course)
	url := apiBaseUrl + "courses/" + courseId + "?apiKey=" + apiKey
	//fmt.Println(url)
	request, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode == http.StatusCreated {
			printResponse(data, "Success")
		} else {
			printResponse(data, "> Error encountered")
		}
	}
}

func deleteCourse(courseId string) {

	url := apiBaseUrl + "courses/" + courseId + "?apiKey=" + apiKey
	//fmt.Println(url)
	request, _ := http.NewRequest(http.MethodDelete, url, nil)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode == http.StatusCreated {
			printResponse(data, "Success")
		} else {
			printResponse(data, "> Error encountered")
		}
	}
}

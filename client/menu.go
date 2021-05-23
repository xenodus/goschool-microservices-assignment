package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var preLoginMenuOptions []string
var postLoginMenuOptions []string

func init() {

	preLoginMenuOptions = []string{
		"Register (Get API Key)",
		"Login (Enter API Key)",
	}

	postLoginMenuOptions = []string{
		"Get a course",
		"Get all courses",
		"Get only active courses",
		"Get only inactive courses",
		"Create course",
		"Update course",
		"Delete course",
		"Show API Key",
		"Invalidate Key (admin login required)",
		"Logout / Change API Key",
	}
}

func menu() {

	for apiKey == "" {
		preLoginMenu()
	}

	postLoginMenu()
}

func preLoginMenu() {
	var preMenuActionNo int

	for apiKey == "" {
		var preMenuActionNoString string

		fmt.Println("Select an option (0 to exit program):")
		for k, v := range preLoginMenuOptions {
			fmt.Println(strconv.Itoa(k+1)+".", v)
		}

		fmt.Scanln(&preMenuActionNoString)

		if preMenuActionNoString == "0" {
			fmt.Println("Exiting program")
			os.Exit(0)
		}

		if i, err := strconv.Atoi(preMenuActionNoString); err == nil {
			preMenuActionNo = i

			if preMenuActionNo > 0 && preMenuActionNo <= len(preLoginMenuOptions) {
				break
			}
		}

		fmt.Println("Error: Invalid Choice!")
	}

	switch preMenuActionNo {
	case 1:
		register() // register
	default:
		for apiKey == "" {
			fmt.Println("Enter API Key to proceed (0 to abort):")
			fmt.Scanln(&apiKey)
		}

		if apiKey == "0" {
			apiKey = ""
			return
		}
	}
}

func postLoginMenu() {
	var actionNo int

	for {
		var actionNoString string

		fmt.Println("Select an option (0 to exit program):")
		for k, v := range postLoginMenuOptions {
			fmt.Println(strconv.Itoa(k+1)+".", v)
		}

		fmt.Scanln(&actionNoString)

		if actionNoString == "0" {
			fmt.Println("Exiting program")
			os.Exit(0)
		}

		if i, err := strconv.Atoi(actionNoString); err == nil {
			actionNo = i

			if actionNo > 0 && actionNo <= len(postLoginMenuOptions) {
				break
			}
		}

		fmt.Println("Error: Invalid Choice!")
	}

	switch actionNo {
	case 1:
		getCoursePrompt() // get a course
	case 2:
		getCourses("all") // get all courses
	case 3:
		getCourses("active") // get only active courses
	case 4:
		getCourses("inactive") // get only inactive courses
	case 5:
		createCourse() // create course
	case 6:
		updateCourse() // update course
	case 7:
		deleteCoursePrompt() // delete course
	case 8:
		fmt.Println("API Key:", apiKey)
	case 9:
		invalidateKey()
	default:
		logout()
	}
}

func courseIDPrompt() string {
	var courseId string = ""

	for courseId == "" {
		fmt.Println("Enter course ID (0 to abort):")
		fmt.Scanln(&courseId)
	}

	return courseId
}

func getCoursePrompt() {

	courseId := courseIDPrompt()

	if courseId != "0" {
		getCourse(courseId)
	}
}

func deleteCoursePrompt() {

	courseId := courseIDPrompt()

	if courseId != "0" {
		deleteCourse(courseId)
	}
}

func newCourseInfoPrompt() *Course {

	courseId, title, description, status := "", "", "", ""

	for courseId == "" {
		fmt.Println("Enter course ID (0 to abort):")
		fmt.Scanln(&courseId)
	}

	if courseId == "0" {
		return nil
	}

	for title == "" {
		fmt.Println("Enter course title (0 to abort):")
		fmt.Scanln(&title)
	}

	if title == "0" {
		return nil
	}

	for description == "" {
		fmt.Println("Enter course description (0 to abort):")
		fmt.Scanln(&description)
	}

	if description == "0" {
		return nil
	}

	for status != "1" && status != "2" && status != "0" {
		fmt.Println("Enter course status (0 to abort):")
		fmt.Println("1. Active")
		fmt.Println("2. Inactive")
		fmt.Scanln(&status)
	}

	if status == "0" {
		return nil
	} else if status == "1" {
		status = "active"
	} else {
		status = "inactive"
	}

	c := Course{courseId, title, description, status}

	return &c
}

func createCourse() {
	c := newCourseInfoPrompt()

	if c != nil {
		c.createCourse()
	}
}

func updateCourseInfoPrompt() (string, *Course) {

	currentCourseId, courseId, title, description, status := "", "", "", "", ""

	for currentCourseId == "" {
		fmt.Println("Enter the course ID you would like to update (0 to abort):")
		fmt.Scanln(&currentCourseId)
	}

	if currentCourseId == "0" {
		return "", nil
	}

	for courseId == "" {
		fmt.Println("Enter the new course ID (0 to abort):")
		fmt.Scanln(&courseId)
	}

	if courseId == "0" {
		return "", nil
	}

	for title == "" {
		fmt.Println("Enter the new course title (0 to abort):")
		fmt.Scanln(&title)
	}

	if title == "0" {
		return "", nil
	}

	for description == "" {
		fmt.Println("Enter the new course description (0 to abort):")
		fmt.Scanln(&description)
	}

	if description == "0" {
		return "", nil
	}

	for status != "1" && status != "2" && status != "0" {
		fmt.Println("Enter the new course status (0 to abort):")
		fmt.Println("1. Active")
		fmt.Println("2. Inactive")
		fmt.Scanln(&status)
	}

	if status == "0" {
		return "", nil
	} else if status == "1" {
		status = "active"
	} else {
		status = "inactive"
	}

	c := Course{courseId, title, description, status}

	return currentCourseId, &c
}

func updateCourse() {
	currentCourseId, c := updateCourseInfoPrompt()

	if c != nil {
		c.updateCourse(currentCourseId)
	}
}

func register() {

	email, password := "", ""

	for email == "" {
		fmt.Println("Enter an email address (0 to abort):")
		fmt.Scanln(&email)

		if email == "0" {
			return
		}

		err := validateEmail(email)

		if err != nil {
			email = ""
			fmt.Println(err.Error())
		}
	}

	for password == "" {
		fmt.Println("Enter a password (0 to abort):")
		inputReader := bufio.NewReader(os.Stdin)
		input, _ := inputReader.ReadString('\n')
		password = strings.TrimRight(input, "\r\n")

		if password == "0" {
			return
		}

		err := validatePassword(password)

		if err != nil {
			password = ""
			fmt.Println(err.Error())
		}
	}

	u := User{email, password}
	u.register()
}

func invalidateKey() {
	email, password, key := "", "", ""

	for email == "" {
		fmt.Println("Enter your email address (0 to abort):")
		fmt.Scanln(&email)

		if email == "0" {
			return
		}

		err := validateEmail(email)

		if err != nil {
			email = ""
			fmt.Println(err.Error())
		}
	}

	for password == "" {
		fmt.Println("Enter your password (0 to abort):")
		inputReader := bufio.NewReader(os.Stdin)
		input, _ := inputReader.ReadString('\n')
		password = strings.TrimRight(input, "\r\n")

		if password == "0" {
			return
		}

		err := validatePassword(password)

		if err != nil {
			password = ""
			fmt.Println(err.Error())
		}
	}

	for key == "" {
		fmt.Println("Enter the API Key to invalidate (0 to abort):")
		fmt.Scanln(&key)

		if key == "0" {
			return
		}
	}

	u := User{email, password}
	u.invalidateKey(key)
}

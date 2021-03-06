package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// drop and create tables
func setupDb() {

	fmt.Println("> resetApp is set to TRUE in conf.go. Proceeding to reset database.")

	// Start DB setup
	fmt.Println("Step 1 of 3: Setting up database")

	// Course
	myDb.Query("DROP TABLE course")
	myDb.Query(`CREATE TABLE course (
		Id varchar(20) NOT NULL PRIMARY KEY,
		Title varchar(255) NOT NULL,
		Description varchar(255) NOT NULL,
		Status ENUM('inactive', 'active') NOT NULL
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	// User
	myDb.Query("DROP TABLE user")
	myDb.Query(`CREATE TABLE user (
		Id varchar(128) NOT NULL PRIMARY KEY,
		Email varchar(128) NOT NULL,
		Password varchar(255) NOT NULL,
		ApiKeyId char(11),
		Admin tinyint(1) NOT NULL DEFAULT '0'
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	myDb.Query("ALTER TABLE `user` ADD UNIQUE( `ApiKeyId`)")
	myDb.Query("ALTER TABLE `user` ADD UNIQUE( `Email`)")

	// API Key
	myDb.Query("DROP TABLE apikey")
	myDb.Query(`CREATE TABLE apikey (
			Id int(11) PRIMARY KEY AUTO_INCREMENT,
			Value varchar(128) NOT NULL,
			Status ENUM('inactive', 'active') NOT NULL
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	myDb.Query("ALTER TABLE `apikey` ADD UNIQUE( `Value`)")

	seedData()

	fmt.Println("Finished resetting app!")
}

// seed test data
func seedData() {
	fmt.Println("Step 2 of 3: Start seeding courses")

	courses := []Course{
		Course{
			"PY101", "Python 101", "Python is a programming language that lets you work more quickly and integrate your systems more effectively.", "inactive",
		},
		Course{
			"GO401", "Golang by GoSchool ", "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.", "active",
		},
		Course{
			"DO321", "Docker for Beginners", "Learn to build and deploy your distributed applications easily to the cloud with Docker.", "active",
		},
	}

	for _, v := range courses {
		v.create()
		fmt.Println("Created course:", v)
	}

	fmt.Println("End seeding courses")

	fmt.Println("Step 3 of 3: Start seeding admin users")
	emails := []string{
		"xenodus@gmail.com",
		"contact@alvinyeoh.com",
		"test@test.com",
	}

	for _, email := range emails {

		userid := uuid.Must(uuid.NewV4()).String()
		bPassword, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.MinCost)
		u := User{userid, email, string(bPassword), -1, 1}

		u.register()
		k, e := u.getKey()

		if e == nil {
			fmt.Println("Created User:", u.Email, "ApiKey:", k.Value)
		}
	}

	fmt.Println("End seeding users")
}

func doLog(logType, msg string) {

	file, err := os.OpenFile("./logs/out.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	defer file.Close()

	var logger *log.Logger

	logType = strings.ToUpper(logType)

	if logType == "INFO" {
		logger = log.New(io.MultiWriter(os.Stdout, file), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else if logType == "WARNING" {
		logger = log.New(io.MultiWriter(os.Stderr, file), "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else if logType == "ERROR" {
		logger = log.New(io.MultiWriter(os.Stderr, file), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logger = log.New(ioutil.Discard, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	logger.Println(msg)
}

func printHeader() {

	const headerText = `
    __  ___ _                                              _
   /  |/  /(_)_____ _____ ____   _____ ___   _____ _   __ (_)_____ ___   _____
  / /|_/ // // ___// ___// __ \ / ___// _ \ / ___/| | / // // ___// _ \ / ___/
 / /  / // // /__ / /   / /_/ /(__  )/  __// /    | |/ // // /__ /  __/(__  )
/_/  /_//_/ \___//_/    \____//____/ \___//_/     |___//_/ \___/ \___//____/

   _____
  / ___/ ___   _____ _   __ ___   _____
  \__ \ / _ \ / ___/| | / // _ \ / ___/
 ___/ //  __// /    | |/ //  __// /
/____/ \___//_/     |___/ \___//_/
`

	fmt.Println(headerText)
}

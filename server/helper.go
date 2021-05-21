package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// drop and create tables
func setupDb() {
	// Start  DB setup
	fmt.Println("Start DB setup")

	// Course
	myDb.Query("DROP TABLE course")
	myDb.Query(`CREATE TABLE course (
		Id varchar(20) NOT NULL PRIMARY KEY,
		Title varchar(255) NOT NULL,
		Description varchar(255) NOT NULL,
		Status ENUM('inactive', 'active') NOT NULL
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	// API Key
	myDb.Query("DROP TABLE apikey")
	myDb.Query(`CREATE TABLE apikey (
		Id int(11) PRIMARY KEY AUTO_INCREMENT,
		Value varchar(128) NOT NULL,
		Status ENUM('inactive', 'active') NOT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	myDb.Query("ALTER TABLE `apikey` ADD UNIQUE( `Value`)")

	// End DB setup
	fmt.Println("Finished DB setup")

	seedData()
}

// seed test data
func seedData() {
	fmt.Println("Start seeding Course")

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
		v.createCourse()
		fmt.Println("Created course:", v)
	}

	fmt.Println("End seeding Course")

	fmt.Println("Start seeding api keys")
	for keys2generate := 0; keys2generate < 2; keys2generate++ {
		k, e := generateKey()

		if e != nil {
			fmt.Println(e.Error())
		} else {
			fmt.Println("Generated key:", k)
		}
	}
	fmt.Println("End seeding api keys")
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

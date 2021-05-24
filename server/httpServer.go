package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// startWebServer is for setting up routes, listening and serving the http server
func startWebServer() {
	router := mux.NewRouter()

	// get all courses
	router.HandleFunc("/api/v1/courses", getCoursesHandler("active"))
	router.HandleFunc("/api/v1/courses/active", getCoursesHandler("active"))
	router.HandleFunc("/api/v1/courses/inactive", getCoursesHandler("inactive"))
	router.HandleFunc("/api/v1/courses/all", getCoursesHandler("all"))
	// get a course
	router.HandleFunc("/api/v1/courses/{courseid}", getCourseHandler).Methods("GET")
	// create a course
	router.HandleFunc("/api/v1/courses/{courseid}", createCourseHandler).Methods("POST")
	// update a course
	router.HandleFunc("/api/v1/courses/{courseid}", updateCourseHandler).Methods("PUT")
	// delete a course
	router.HandleFunc("/api/v1/courses/{courseid}", deleteCourseHandler).Methods("DELETE")

	// user
	router.HandleFunc("/api/v1/users", registerHandler).Methods("POST")
	// invalidate an api key
	router.HandleFunc("/api/v1/keys/{apiKey}", invalidateKeyHandler).Methods("DELETE")

	fmt.Println("Listening at http://" + serverHostname + ":" + serverPort)

	err := http.ListenAndServe(serverHostname+":"+serverPort, router)
	if err != nil {
		log.Fatal(err)
	}
}

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

	// Routes

	// Get all courses
	router.HandleFunc("/api/v1/courses", getCoursesHandler("active"))
	router.HandleFunc("/api/v1/courses/active", getCoursesHandler("active"))
	router.HandleFunc("/api/v1/courses/inactive", getCoursesHandler("inactive"))
	router.HandleFunc("/api/v1/courses/all", getCoursesHandler("all"))
	// Get a course
	router.HandleFunc("/api/v1/courses/{courseid}", getCourseHandler).Methods("GET")
	// Create a course
	router.HandleFunc("/api/v1/courses/{courseid}", createCourseHandler).Methods("POST")
	// Update a course
	router.HandleFunc("/api/v1/courses/{courseid}", updateCourseHandler).Methods("PUT")
	// Delete a course
	router.HandleFunc("/api/v1/courses/{courseid}", deleteCourseHandler).Methods("DELETE")

	// User
	router.HandleFunc("/api/v1/users", registerHandler).Methods("POST")
	// Invalidate an api key
	router.HandleFunc("/api/v1/keys/{apiKey}", invalidateKeyHandler).Methods("DELETE")

	fmt.Println("Listening at port 8080")

	err := http.ListenAndServe(serverHostname+":"+serverPort, router)
	if err != nil {
		log.Fatal(err)
	}
}

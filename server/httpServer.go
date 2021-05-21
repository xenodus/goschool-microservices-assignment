package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func testHandler(res http.ResponseWriter, req *http.Request) {

	c, e1 := getCourse("PY101")
	fmt.Println(c, e1)

	d, e2 := getCourse("1312414")
	fmt.Println(d, e2)

	c.Title = "lalala"
	c.updateCourse()

	e, _ := getCourse("PY101")
	fmt.Println(e)

	fmt.Println("Yay!")
}

// apiCreateCourse is for creating courses via POST requests
func apiCreateCourse(res http.ResponseWriter, req *http.Request) {

	if !isKeyValid(req) {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusUnauthorized, errInvalidApiKey.Error()})
		return
	}

	if req.Header.Get("Content-type") == "application/json" {
		var newCourse Course
		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {

			// convert JSON to object
			json.Unmarshal(reqBody, &newCourse)

			// check for invalid values and db lengths
			validateErr := newCourse.validateFields()

			if validateErr != nil {
				res.Header().Set("Content-Type", "application/json")
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusUnprocessableEntity, validateErr.Error()})
				return
			}

			// check if course exists; add only if course does not exist
			// course does not exist
			_, cErr := getCourse(newCourse.Id)

			// course exists
			if cErr == nil {
				res.Header().Set("Content-Type", "application/json")
				res.WriteHeader(http.StatusConflict)
				json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusConflict, errDuplicateCourseId.Error()})
				return
			}

			// create course
			if cErr == errCourseNotFound {
				doLog("INFO", req.RemoteAddr+" | Created course: "+newCourse.Id)
				newCourse.createCourse()

				res.Header().Set("Content-Type", "application/json")
				res.WriteHeader(http.StatusCreated)
				json.NewEncoder(res).Encode(JSONResponse{"ok", http.StatusCreated, "Course created, " + newCourse.Id})
				return
			} else {
				doLog("ERROR", err.Error())

				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(errInternalServerError.Error()))
				return
			}
		}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusUnprocessableEntity, errInvalidCourseInfo.Error()})
}

// apiUpdateCourse is for creating OR editing courses via PUT requests
func apiUpdateCourse(res http.ResponseWriter, req *http.Request) {

	if !isKeyValid(req) {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusUnauthorized, errInvalidApiKey.Error()})
		return
	}

	if req.Header.Get("Content-type") == "application/json" {
		var newCourse Course
		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {
			// convert JSON to object
			json.Unmarshal(reqBody, &newCourse)

			// check for invalid values and db lengths
			validateErr := newCourse.validateFields()

			if validateErr != nil {
				res.Header().Set("Content-Type", "application/json")
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusUnprocessableEntity, validateErr.Error()})
				return
			}
		}

		params := mux.Vars(req)

		if id, ok := params["courseid"]; ok {

			course, err := getCourse(id)

			if err != nil {
				// create
				if err == errCourseNotFound {
					newCourse.createCourse()
					doLog("INFO", req.RemoteAddr+" | Created course: "+newCourse.Id)

					res.Header().Set("Content-Type", "application/json")
					res.WriteHeader(http.StatusCreated)
					json.NewEncoder(res).Encode(JSONResponse{"ok", http.StatusCreated, "Course created, " + newCourse.Id})
					return
				} else {
					doLog("ERROR", err.Error())
				}
			}

			// update
			if course != nil {
				course = &newCourse
				course.updateCourse()
				doLog("INFO", req.RemoteAddr+" | Updated course: "+course.Id)

				res.Header().Set("Content-Type", "application/json")
				json.NewEncoder(res).Encode(JSONResponse{"ok", http.StatusOK, "Course updated, " + course.Id})
				return
			}
		}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusUnprocessableEntity, errInvalidCourseInfo.Error()})
}

// apiDeleteCourse is for deleting courses via DELETE requests
func apiDeleteCourse(res http.ResponseWriter, req *http.Request) {

	if !isKeyValid(req) {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusUnauthorized, errInvalidApiKey.Error()})
		return
	}

	params := mux.Vars(req)

	if id, ok := params["courseid"]; ok {

		course, err := getCourse(id)

		if err != nil {
			if err == errCourseNotFound {
				doLog("WARNING", err.Error())
			} else {
				doLog("ERROR", err.Error())
			}
		}

		// delete
		if course != nil {
			doLog("INFO", req.RemoteAddr+" | Deleted course: "+course.Id)
			course.deleteCourse()

			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(JSONResponse{"ok", http.StatusOK, "Course deleted"})
			return
		}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusNotFound)
	json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusNotFound, errCourseNotFound.Error()})
}

// apiGetCourse is for getting a course via GET requests
func apiGetCourse(res http.ResponseWriter, req *http.Request) {

	if !isKeyValid(req) {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusUnauthorized, errInvalidApiKey.Error()})
		return
	}

	params := mux.Vars(req)

	if id, ok := params["courseid"]; ok {

		course, err := getCourse(id)

		if err != nil {
			if err == errCourseNotFound {
				doLog("WARNING", err.Error())
			} else {
				doLog("ERROR", err.Error())
			}
		}

		if course != nil {
			json.NewEncoder(res).Encode(course)
			return
		}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusNotFound)
	json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusNotFound, errCourseNotFound.Error()})
}

// apiGetCourses is for getting courses via GET requests
func apiGetCourses(status string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if !isKeyValid(req) {
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusUnauthorized, errInvalidApiKey.Error()})
			return
		}

		courses, err := getCourses(status)

		if err != nil {
			doLog("ERROR", err.Error())
		}

		// Includes empty result
		json.NewEncoder(res).Encode(courses)
	}
}

func apiInvalidateKey(res http.ResponseWriter, req *http.Request) {

	key := req.FormValue("adminKey")

	if key != adminKey {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusUnauthorized, errInvalidApiKey.Error()})
		return
	}

	params := mux.Vars(req)

	if apiKey, ok := params["apiKey"]; ok {
		err := invalidateKey(apiKey)

		if err == nil {
			doLog("INFO", req.RemoteAddr+" | Invalidated api key: "+apiKey)

			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(JSONResponse{"ok", http.StatusOK, "Api key invalidated"})
			return
		}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusNotFound)
	json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusNotFound, errCourseNotFound.Error()})
}

// startWebServer is for setting up routes, listening and serving the http server
func startWebServer() {
	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/api/v1/test", testHandler)

	// Get all courses
	router.HandleFunc("/api/v1/courses", apiGetCourses("active"))
	router.HandleFunc("/api/v1/courses/active", apiGetCourses("active"))
	router.HandleFunc("/api/v1/courses/inactive", apiGetCourses("inactive"))
	router.HandleFunc("/api/v1/courses/all", apiGetCourses("all"))
	// Get a course
	router.HandleFunc("/api/v1/courses/{courseid}", apiGetCourse).Methods("GET")
	// Create a course
	router.HandleFunc("/api/v1/courses/{courseid}", apiCreateCourse).Methods("POST")
	// Update a course
	router.HandleFunc("/api/v1/courses/{courseid}", apiUpdateCourse).Methods("PUT")
	// Delete a course
	router.HandleFunc("/api/v1/courses/{courseid}", apiDeleteCourse).Methods("DELETE")

	// Invalidating an api key - not secured
	router.HandleFunc("/api/v1/key/{apiKey}", apiInvalidateKey).Methods("DELETE")

	fmt.Println("Listening at port 8080")

	err := http.ListenAndServe(serverHostname+":"+serverPort, router)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// createCourseHandler is for creating courses via POST requests
func createCourseHandler(res http.ResponseWriter, req *http.Request) {

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
				newCourse.create()

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

// updateCourseHandler is for creating OR editing courses via PUT requests
func updateCourseHandler(res http.ResponseWriter, req *http.Request) {

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
					// use courseid in params instead of json when creating
					newCourse.Id = id
					newCourse.create()
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

				// change in course id so check if new course id in json already exists
				if newCourse.Id != params["courseid"] {
					_, err := getCourse(newCourse.Id)

					// exists
					if err == nil {
						res.Header().Set("Content-Type", "application/json")
						res.WriteHeader(http.StatusConflict)
						json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusConflict, errDuplicateCourseId.Error()})
						return
					}
				}

				course.update(&newCourse)
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

// deleteCourseHandler is for deleting courses via DELETE requests
func deleteCourseHandler(res http.ResponseWriter, req *http.Request) {

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
			course.delete()

			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(JSONResponse{"ok", http.StatusOK, "Course deleted"})
			return
		}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusNotFound)
	json.NewEncoder(res).Encode(JSONResponse{"error", http.StatusNotFound, errCourseNotFound.Error()})
}

// getCourseHandler is for getting a course via GET requests
func getCourseHandler(res http.ResponseWriter, req *http.Request) {

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

// getCoursesHandler is for getting all courses by status type via GET requests
func getCoursesHandler(status string) http.HandlerFunc {
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

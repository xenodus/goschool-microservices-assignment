package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func registerHandler(res http.ResponseWriter, req *http.Request) {

	if req.Header.Get("Content-type") == "application/json" {
		var userRegistration UserAuth
		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {

			// convert JSON to object
			json.Unmarshal(reqBody, &userRegistration)

			// check for invalid values and db lengths
			validateErr := userRegistration.validateFields()

			if validateErr != nil {
				printJSONResponse(res, JSONResponse{"error", http.StatusUnprocessableEntity, validateErr.Error()})
				return
			}

			// check if user exists; add only if user does not exist
			_, uErr := getUserByEmail(userRegistration.Email)

			// user (email) exists
			if uErr == nil {
				printJSONResponse(res, JSONResponse{"error", http.StatusConflict, errDuplicateUser.Error()})
				return
			}

			// create user
			if uErr == errUserNotFound {

				userid := uuid.Must(uuid.NewV4()).String()
				bPassword, _ := bcrypt.GenerateFromPassword([]byte(userRegistration.Password), bcrypt.MinCost)

				u := User{userid, userRegistration.Email, string(bPassword), -1, 0}
				u.register()
				k, _ := u.getKey()
				doLog("INFO", req.RemoteAddr+" | Created user: "+u.Id)

				res.Header().Set("Content-Type", "application/json")
				res.WriteHeader(http.StatusCreated)
				json.NewEncoder(res).Encode(ApiKeyResponse{"ok", http.StatusCreated, k.Value, k.Status})
				return
			} else {
				doLog("ERROR", err.Error())

				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(errInternalServerError.Error()))
				return
			}
		}
	}

	printJSONResponse(res, JSONResponse{"error", http.StatusUnprocessableEntity, errInvalidUserInfo.Error()})
}

func invalidateKeyHandler(res http.ResponseWriter, req *http.Request) {

	if !isKeyValid(req) {
		handleKeyInvalid(res)
		return
	}

	if req.Header.Get("Content-type") == "application/json" {
		var authUser UserAuth
		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {

			// convert JSON to object
			json.Unmarshal(reqBody, &authUser)

			// check for invalid values and db lengths
			validateErr := authUser.validateFields()

			if validateErr != nil {
				printJSONResponse(res, JSONResponse{"error", http.StatusUnprocessableEntity, validateErr.Error()})
				return
			}

			// check if user exists; add only if user does not exist
			user, uErr := getUserByEmail(authUser.Email)

			// user don't exists
			if uErr != nil {
				printJSONResponse(res, JSONResponse{"error", http.StatusUnauthorized, errAuthFailure.Error()})
				return
			}

			// check if password matches
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authUser.Password))

			if err != nil {
				printJSONResponse(res, JSONResponse{"error", http.StatusUnauthorized, errAuthFailure.Error()})
				return
			}

			// check if user is an admin
			if user.Admin == 0 {
				printJSONResponse(res, JSONResponse{"error", http.StatusForbidden, errNoPermission.Error()})
				return
			}

			params := mux.Vars(req)

			if apiKey, ok := params["apiKey"]; ok {

				err := invalidateKey(apiKey)

				if err == nil {
					doLog("INFO", req.RemoteAddr+" | Invalidated api key: "+apiKey)
					printJSONResponse(res, JSONResponse{"ok", http.StatusOK, "api key invalidated"})
					return
				}
			}

			printJSONResponse(res, JSONResponse{"error", http.StatusUnprocessableEntity, errKeyInvalidateNotFound.Error()})
			return
		}
	}

	printJSONResponse(res, JSONResponse{"error", http.StatusUnauthorized, errAuthFailure.Error()})
}

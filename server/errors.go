package main

import (
	"errors"
)

// Errors and accompanying messages to be output in logs or to users.
var (
	// User
	errUserNotFound    = errors.New("user not found")
	errInvalidUserInfo = errors.New("please supply user information in JSON format")
	errDuplicateUser   = errors.New("the email address is taken")
	errAuthFailure     = errors.New("invalid login information")
	errNoPermission    = errors.New("insufficient rights")

	// Api Key
	errKeyNotFound           = errors.New("api key not found")
	errKeyInvalidateNotFound = errors.New("api key to invalidate not found")

	// Course
	errCourseNotFound      = errors.New("course not found")
	errDuplicateCourseId   = errors.New("duplicate course id")
	errInvalidCourseInfo   = errors.New("please supply course information in JSON format")
	errInvalidCourseStatus = errors.New("invalid course status. Please supply either active or inactive for course status")

	// DB Lengths
	errInvalidCourseIdLength          = errors.New("invalid course id length. Please supply course id that's less than 20 characters")
	errInvalidCourseTitleLength       = errors.New("invalid course title length. Please supply course title that's less than 255 characters")
	errInvalidCourseDescriptionLength = errors.New("invalid course description length. Please supply course description that's less than 255 characters")
	errInvalidEmailLength             = errors.New("invalid email length. Please supply course id that's less than 128 characters")

	errInvalidApiKey = errors.New("your API key is invalid or incorrect")

	errInternalServerError = errors.New("an unexpected error has occured")
)

package main

import (
	"errors"
	"strconv"
)

// Errors and accompanying messages to be output in logs or to users.
var (
	// Registration
	errAuthEmailLength    = errors.New("invalid email entered. Please ensure email is valid and between " + strconv.Itoa(emailMinLen) + " and " + strconv.Itoa(emailMaxLen) + " in length")
	errAuthEmailFormat    = errors.New("invalid email entered")
	errAuthPasswordLength = errors.New("invalid password entered. Please ensure password is between " + strconv.Itoa(passwordMinLen) + " and " + strconv.Itoa(passwordMaxLen) + " in length")
)

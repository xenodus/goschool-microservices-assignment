package main

type JSONResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

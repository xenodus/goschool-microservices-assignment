package main

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ApiKeyResponse struct {
	Status       string `json:"status"`
	Code         int    `json:"code"`
	ApiKey       string `json:"apiKey"`
	ApiKeyStatus string `json:"apiKeyStatus"`
}

func printJSONResponse(res http.ResponseWriter, r JSONResponse) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(r.Code)
	json.NewEncoder(res).Encode(r)
}

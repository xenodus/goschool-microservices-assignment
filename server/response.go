package main

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

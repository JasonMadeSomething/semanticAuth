package handlers

import (
	"encoding/json"
	"net/http"
)

// StandardResponse is a consistent response structure for all API endpoints
type StandardResponse struct {
	Status  string      `json:"status"`  // "success" or "error"
	Success bool        `json:"success"` // true or false
	Message string      `json:"message,omitempty"` // optional message
	Data    interface{} `json:"data,omitempty"`    // payload data
}

// RespondWithSuccess sends a standardized success response with data
func RespondWithSuccess(w http.ResponseWriter, message string, data interface{}) {
	response := StandardResponse{
		Status:  "success",
		Success: true,
		Message: message,
		Data:    data,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RespondWithError sends a standardized error response
func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	response := StandardResponse{
		Status:  "error",
		Success: false,
		Message: message,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

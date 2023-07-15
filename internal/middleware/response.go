package middleware

import (
	"net/http"
	"encoding/json"
)

var bodyErrorResponse = make(map[string]string)

func NewJSONResponse(w http.ResponseWriter, data interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
        
		NewJSONResponseWhithError(
			w, 
			http.StatusInternalServerError, 
			"data acquisition error...",
		)
	} 
}

func NewJSONResponseWhithError(
	w http.ResponseWriter, statusCode int, message string,
) {

	w.WriteHeader(statusCode)

	if message != "" {
		bodyErrorResponse["error"] = message
		json.NewEncoder(w).Encode(&bodyErrorResponse)
	}
}

package util

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

func RespondError(w http.ResponseWriter, msg string, code int) {
	errorResponse := Error{Message: msg}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorResponse)
}

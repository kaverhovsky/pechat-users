package http_helpers

import (
	"encoding/json"
	"net/http"
)

// TODO move to pechat lib

func WriteSuccess(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			WriteError(w, http.StatusInternalServerError, err)
		}
	}
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if jsonErr := json.NewEncoder(w).Encode(err); jsonErr != nil {
		panic(err)
	}
}

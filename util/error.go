package util

import (
	"encoding/json"
	"net/http"

	"github.com/diontristen/go-crud/models"
)

type ErrorJSON struct {
	Error models.ValidationResponse `json:"error"`
}

// This utility is used for specific Error Returns.

// ReturnBadRequest is a helper that is meant to be used when some error is caught.
// You can pass a message to the user, or let it return the error formated as string.
func ReturnBadRequest(w http.ResponseWriter, err string) {
	returnError := models.DefaultError{Error: err}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(returnError)
}

func ReturnBadRequestJSON(w http.ResponseWriter, err models.ValidationResponse) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorJSON{Error: err})
}

// ReturnUnauthorized is a helper that is meant to be used when some error is caught.
// You can pass a message to the user, or let it return the error formated as string.
func ReturnUnauthorized(w http.ResponseWriter, err string) {
	returnError := models.DefaultError{Error: err}

	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(returnError)
}

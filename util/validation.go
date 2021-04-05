package util

import (
	"fmt"
	"net/http"

	"github.com/diontristen/go-crud/models"
	"gopkg.in/go-playground/validator.v9"
)

// USAGE: Check if the array s contains the element e
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ValidateData takes a struct with notations and validates it using validator library.
// If there are errors, returns an array with the errors
func ValidateData(d interface{}) (models.ValidationResponse, error) {
	v := validator.New()

	validationErr := v.Struct(d)

	if validationErr != nil {
		lst := models.ValidationResponse{}
		var s models.ValidationError
		for _, e := range validationErr.(validator.ValidationErrors) {
			s = models.ValidationError{
				Key:   e.Field(),
				Value: e.ActualTag(),
			}
			lst = append(lst, s)

		}
		return lst, nil
	}
	return nil, nil
}

func ValidateExisting(w http.ResponseWriter, key string, value interface{}, ac *AppContext, method string) bool {
	var c int

	checkSQL := fmt.Sprintf("SELECT count(*) c FROM contacts WHERE %s = '%v'", key, value)

	checkError := ac.DB.QueryRow(checkSQL).Scan(&c)

	if checkError != nil {
		ReturnBadRequest(w, fmt.Sprint(checkError))
		return false
	}

	if c > 0 && method == "post" {
		ReturnBadRequest(w, "Contact already exists.")
		return false
	}

	if c <= 0 && method == "update" {
		ReturnBadRequest(w, "Contact doesn't exists.")
		return false
	}

	return true
}

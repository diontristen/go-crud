package util

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/diontristen/go-crud/models"
)

// This utility is used for automatic returns base on the error type.

type ErrNotFound interface {
	NotFound()
}

type ObjNotFoundError struct {
	Object string
	ID     string
}

func (e ObjNotFoundError) NotFound() {}

func (e ObjNotFoundError) Error() string {
	return "No such " + e.Object + ": " + e.ID
}

type ErrInvalidParameter interface {
	InvalidParameter()
}

type errInvalidParameter struct {
	error
}

func (errInvalidParameter) InvalidParameter() {}

func InvalidParameter(err error) error {
	if err == nil {
		return err
	}
	return errInvalidParameter{err}
}

type ErrUnauthorized interface {
	Unauthorized()
}

type errUnauthorized struct {
	error
}

func (errUnauthorized) Unauthorized() {}

func (e errUnauthorized) Cause() error {
	return e.error
}

func Unauthorized(err error) error {
	if err == nil {
		return err
	}

	return errUnauthorized{err}
}

type ErrForbidden interface {
	Forbidden()
}

type errForbidden struct{ error }

func (errForbidden) Forbidden() {}

func (e errForbidden) Cause() error {
	return e.error
}

func Forbidden(err error) error {
	if err == nil {
		return err
	}
	return errForbidden{err}
}

// USAGE: JSON Encode to be written/returned as a response to the request with the status code 200
func WriteJSONWithStatus(w http.ResponseWriter, v interface{}, statusCode int) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		Errorf("json.Marshal: %v", err)
		WriteHTTPError(w, err)
	}

	w.WriteHeader(statusCode)
	// nolint: errcheck
	w.Write(jsonBytes)
}

// USAGE: Calls the WriteJSONWithStatus
func WriteJSON(w http.ResponseWriter, v interface{}) {
	WriteJSONWithStatus(w, v, http.StatusOK)
}

// USAGE: JSON Encode Error Responses.
func serializeError(w http.ResponseWriter, errStr string) {
	err := json.NewEncoder(w).Encode(models.DefaultError{Error: errStr})
	if err != nil {
		Error(err)
	}
}

// USAGE: Checks what Status Code should be returned
func WriteHTTPError(w http.ResponseWriter, respErr error) {
	if respErr == nil {
		err := errors.New("unexpected HTTP error handling")
		Error(err)
		respErr = err
	}

	statusCode := http.StatusInternalServerError // Response Code 500

	// You can check status codes here: https://golang.org/pkg/net/http/
	switch respErr.(type) {
	case ErrNotFound:
		statusCode = http.StatusNotFound
	case ErrInvalidParameter:
		statusCode = http.StatusBadRequest
	case ErrUnauthorized:
		statusCode = http.StatusUnauthorized
	case ErrForbidden:
		statusCode = http.StatusForbidden
	}

	w.WriteHeader(statusCode)

	serializeError(w, respErr.Error())

}

// USAGE: Handler of Error Response
func WriteError(w http.ResponseWriter, err error) {
	Error(err)
	WriteHTTPError(w, err)
}

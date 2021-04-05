package models

// DefaultError is the stantdard interface to return errors to the client's api
type DefaultError struct {
	Error string `json:"error"`
}

type ValidationError struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ValidationResponse []ValidationError

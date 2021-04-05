package util

import (
	"encoding/json"
	"fmt"
)

// USAGE: Adds indentation to the JSON Data
func MarshalIndent(v interface{}) ([]byte, error) {
	dat, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return dat, fmt.Errorf("Error marshalling JSON: %s", err)
	}

	return dat, nil
}

// USAGE: JSON Indentation Handler
func DumpIndent(v interface{}) {
	b, err := MarshalIndent(v)
	if err != nil {
		Errorf("MarshalIndent: %v", err)
	}

	fmt.Println(string(b))
}

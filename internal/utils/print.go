package utils

import (
	"encoding/json"
	"fmt"
)

// PrettyPrint prints the given interface in a pretty format
func PrettyPrint(v ...interface{}) {
	for _, i := range v {
		b, err := json.MarshalIndent(i, "", "  ")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	}
}

// PrettyString returns the given interface in a pretty format
func PrettyString(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// AnyToString converts any type to string
func AnyToString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

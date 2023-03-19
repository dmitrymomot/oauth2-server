package utils

import "fmt"

// GetVarType returns the type of the given variable as a string.
func GetVarType(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

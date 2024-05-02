package converter

import (
	"encoding/json"
)

// CommaStringArray: Convert string array to comma separated string.
// Example: CommaStringArray([]string{"a", "b", "c"}) => `["a","b","c"]`
func CommaStringArray(data []string) string {
	jsonBytes, _ := json.Marshal(data)
	return string(jsonBytes)
}

// StructToString: Convert struct to string.
// Example: StructToString(testStruct{"a": "b"}) => `{"a": "b"}`
func StructToString[T any](data T) string {
	jsonBytes, _ := json.Marshal(data)
	return string(jsonBytes)
}

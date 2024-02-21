package converter

import (
	"encoding/json"
	"strings"
)

// CommaStringArray: Convert string array to comma separated string.
// Example: CommaStringArray([]string{"a", "b", "c"}) => `["a","b","c"]`
func CommaStringArray(data []string) string {
	var result string
	if len(data) > 0 {
		result = "[\"" + strings.Join(data, "\",\"") + "\"]"
	}
	return result
}

// StructToString: Convert struct to string.
// Example: StructToString(testStruct{"a": "b"}) => `{"a": "b"}`
func StructToString[T any](data T) string {
	jsonBytes, _ := json.Marshal(data)
	return string(jsonBytes)
}

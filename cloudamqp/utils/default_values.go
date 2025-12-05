package utils

func DefaultInt64Value(value, defaultValue int64) int64 {
	if value == 0 {
		return defaultValue
	}
	return value
}

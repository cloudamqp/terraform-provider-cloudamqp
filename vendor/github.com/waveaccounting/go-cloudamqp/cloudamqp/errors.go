package cloudamqp

import "fmt"

// APIError represents a CloudAMQP API Error response
type APIError map[string]interface{}

// TODO: use this instead
// type apiError struct {
// 	Error string `json:"error"`
// }

func (e APIError) Error() string {
	if len(e) == 1 {
		if error, ok := e["error"].(string); ok {
			return fmt.Sprintf("cloudamqp: %s", error)
		}
	}

	return fmt.Sprintf("sentry: %v", map[string]interface{}(e))
}

// Empty returns true if empty.
func (e APIError) Empty() bool {
	return len(e) == 0
}

func relevantError(httpError error, apiError APIError) error {
	if httpError != nil {
		return httpError
	}
	if !apiError.Empty() {
		return apiError
	}
	return nil
}

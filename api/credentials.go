package api

import (
	"fmt"
	"regexp"
	"strconv"
)

func (api *API) ReadCredentials(id int) (map[string]any, error) {
	var (
		data       map[string]any
		failed     map[string]any
		instanceID = strconv.Itoa(id)
	)

	response, err := api.sling.New().Path("/api/instances/").Get(instanceID).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return extractInfo(data["url"].(string)), nil
	default:
		return nil, fmt.Errorf("read credentials failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func extractInfo(url string) map[string]any {
	paramsMap := make(map[string]any)
	r := regexp.MustCompile(`^.*:\/\/(?P<username>(.*)):(?P<password>(.*))@`)
	match := r.FindStringSubmatch(url)

	for i, name := range r.SubexpNames() {
		if name == "username" {
			paramsMap["username"] = match[i]
		}
		if name == "password" {
			paramsMap["password"] = match[i]
		}
	}

	return paramsMap
}

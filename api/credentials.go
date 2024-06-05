package api

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
)

func (api *API) ReadCredentials(id int) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	instanceID := strconv.Itoa(id)
	log.Printf("[DEBUG] go-api::credentials::read instance ID: %v", instanceID)
	response, err := api.sling.New().Path("/api/instances/").Get(instanceID).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadCredentials failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return extractInfo(data["url"].(string)), nil
}

func extractInfo(url string) map[string]interface{} {
	paramsMap := make(map[string]interface{})
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

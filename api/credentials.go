package api

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

func (api *API) ReadCredentials(id int) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	instance_id := strconv.Itoa(id)
	response, err := api.sling.Path("/api/instances/").Get(instance_id).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadCredentials failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return ExtractInfo(data["url"].(string)), nil
}

func ExtractInfo(url string) map[string]interface{} {
	paramsMap := make(map[string]interface{})
	r := regexp.MustCompile(`amqp:\/\/(?P<username>(.*?)):(?P<password>(.*?))@`)
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

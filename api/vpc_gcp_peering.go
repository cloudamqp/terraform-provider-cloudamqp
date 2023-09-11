package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func (api *API) waitForGcpPeeringStatus(instanceID int, peerID string) error {
	for {
		time.Sleep(10 * time.Second)
		data, err := api.ReadVpcGcpPeering(instanceID, peerID)
		if err != nil {
			return err
		}
		rows := data["rows"].([]interface{})
		if len(rows) > 0 {
			for _, row := range rows {
				tempRow := row.(map[string]interface{})
				if tempRow["name"] != peerID {
					continue
				}
				if tempRow["state"] == "ACTIVE" {
					return nil
				}
			}
		}
	}
}

func (api *API) RequestVpcGcpPeering(instanceID int, params map[string]interface{},
	waitOnStatus bool) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("api/instances/%v/vpc-peering", instanceID)
	)

	log.Printf("[DEBUG] go-api::vpc_gcp_peering::request params: %v", params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("request VPC peering failed, status: %v, message: %s", response.StatusCode, failed)
	}

	if waitOnStatus {
		log.Printf("[DEBUG] go-api::vpc_gcp_peering::request waiting for active state")
		api.waitForGcpPeeringStatus(instanceID, data["peering"].(string))
	}
	return data, nil
}

func (api *API) ReadVpcGcpPeering(instanceID int, peerID string) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%v/vpc-peering", instanceID)
	)

	log.Printf("[DEBUG] go-api::vpc_gcp_peering::read instance_id: %v, peer_id: %v", instanceID, peerID)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc_gcp_peering::read data: %v", data)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadRequest failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, nil
}

func (api *API) UpdateVpcGcpPeering(instanceID int, peerID string) (map[string]interface{}, error) {
	return api.ReadVpcGcpPeering(instanceID, peerID)
}

func (api *API) RemoveVpcGcpPeering(instanceID int, peerID string) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%v/vpc-peering/%v", instanceID, peerID)
	)

	log.Printf("[DEBUG] go-api::vpc_gcp_peering::remove instance id: %v, peering id: %v", instanceID, peerID)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return fmt.Errorf("RemoveVpcPeering failed, status: %v, message: %s", response.StatusCode, failed)
	}
	return nil
}

func (api *API) ReadVpcGcpInfo(instanceID int) (map[string]interface{}, error) {
	// Initiale values, 5 attempts and 20 second sleep
	return api.readVpcGcpInfoWithRetry(instanceID, 5, 20)
}

func (api *API) readVpcGcpInfoWithRetry(instanceID, attempts, sleep int) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%v/vpc-peering/info", instanceID)
	)

	log.Printf("[DEBUG] go-api::vpc_gcp_peering::info instance id: %v", instanceID)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc_gcp_peering::info data: %v", data)
	if err != nil {
		return nil, err
	}

	statusCode := response.StatusCode
	log.Printf("[DEBUG] go-api::vpc_gcp_peering::info statusCode: %d", statusCode)
	switch {
	case statusCode == 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			if attempts--; attempts > 0 {
				log.Printf("[INFO] go-api::vpc_gcp_peering::info Timeout talking to backend "+
					"attempts left %d and retry in %d seconds", attempts, sleep)
				time.Sleep(time.Duration(sleep) * time.Second)
				return api.readVpcGcpInfoWithRetry(instanceID, attempts, 2*sleep)
			} else {
				return nil, fmt.Errorf("ReadInfo failed, status: %v, message: %s", response.StatusCode, failed)
			}
		}
	}
	return data, nil
}

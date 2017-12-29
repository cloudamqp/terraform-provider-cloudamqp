package cloudamqp

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceService_List(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/instances", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
      {
        "id": 1234,
        "plan": "lemur",
        "region": "amazon-web-services::us-east-1",
        "name": "test-instance-1"
      },
			{
        "id": 1235,
        "plan": "bunny",
        "region": "amazon-web-services::us-east-1",
        "name": "test-instance-2"
      }
    ]`)
	})

	client := NewClient(httpClient, nil, "")
	instances, _, err := client.Instances.List()
	assert.NoError(t, err)
	expected := []Instance{
		{
			ID:     1234,
			Plan:   "lemur",
			Region: "amazon-web-services::us-east-1",
			Name:   "test-instance-1",
		},
		{
			ID:     1235,
			Plan:   "bunny",
			Region: "amazon-web-services::us-east-1",
			Name:   "test-instance-2",
		},
	}
	assert.Equal(t, expected, instances)
}

func TestInstanceService_Get(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/instances/1234", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
  		"id": 1234,
  		"name": "test-instance-1",
  		"plan": "lemur",
  		"region": "amazon-web-services::us-east-1",
  		"url": "amqp://username:password@jolly-wombat.rmq.cloudamqp.com/abcdefg",
  		"apikey": "3d5fbd52-dc07-4ae3-976f-27bf9604e00b"
		}`)
	})

	client := NewClient(httpClient, nil, "")
	instance, _, err := client.Instances.Get(1234)
	assert.NoError(t, err)
	expected := &Instance{
		ID:     1234,
		Name:   "test-instance-1",
		Plan:   "lemur",
		Region: "amazon-web-services::us-east-1",
		URL:    "amqp://username:password@jolly-wombat.rmq.cloudamqp.com/abcdefg",
		APIKey: "3d5fbd52-dc07-4ae3-976f-27bf9604e00b",
	}
	assert.Equal(t, expected, instance)
}

func TestInstanceService_Create(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/instances", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		// TODO: Figure out how to get this function working
		// assertPostForm(t, map[string]interface{}{
		// 	"name":   "test-instance-new",
		// 	"plan":   "tiger",
		// 	"region": "amazon-web-services::us-east-1",
		// }, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": 1236,
			"message": "Instance created",
			"url": "amqp://username:password@happy-fish.rmq.cloudamqp.com/abcdefg",
			"apikey": "d6e6f799-d6ec-407a-a6e7-925414012121"
		}`)
	})

	client := NewClient(httpClient, nil, "")
	params := &CreateInstanceParams{
		Name:   "test-instance-new",
		Plan:   "lemur",
		Region: "amazon-web-services::us-east-1",
	}
	instance, _, err := client.Instances.Create(params)
	assert.NoError(t, err)

	expected := &Instance{
		ID:     1236,
		Name:   "",
		URL:    "amqp://username:password@happy-fish.rmq.cloudamqp.com/abcdefg",
		APIKey: "d6e6f799-d6ec-407a-a6e7-925414012121",
	}
	assert.Equal(t, expected, instance)
}

func TestInstanceService_Update(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/instances/1236", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "PUT", r)
		// TODO: Figure out how to get this function working
		// assertPostForm(t, map[string]interface{}{
		// 		"name":   "test-instance-renamed",
		// 		"plan":   "rabbit",
		// 		"nodes": 2,
		// 	},
		// }, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"message": "Changes saved"
		}`)
	})

	client := NewClient(httpClient, nil, "")
	params := &UpdateInstanceParams{
		Name:  "test-instance-renamed",
		Plan:  "tiger",
		Nodes: 2,
	}
	_, _, err := client.Instances.Update(1236, params)
	assert.NoError(t, err)
}

func TestInstanceService_Delete(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/instances/1236", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "DELETE", r)
	})

	client := NewClient(httpClient, nil, "")
	_, err := client.Instances.Delete(1236)
	assert.NoError(t, err)
}

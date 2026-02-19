package api

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) waitUntilReady(ctx context.Context, instanceID string) (map[string]any, error) {
	path := fmt.Sprintf("/api/instances/%s", instanceID)
	ctxTimeout, cancel := context.WithTimeout(ctx, 1800*time.Second) // 30 minutes
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("waiting for instance to be ready, instanceID=%s", instanceID))
	attempt := 1

	for {
		if ctxTimeout.Err() != nil {
			return nil, fmt.Errorf("timeout reached while waiting for instance to be ready")
		}

		var (
			data   map[string]any
			failed map[string]any
		)

		tflog.Debug(ctx, fmt.Sprintf("Checking instance ready status, attempt=%d", attempt))
		err := api.callWithRetry(ctxTimeout, api.sling.New().Get(path), retryRequest{
			functionName: "waitUntilReady",
			resourceName: "Instance",
			attempt:      attempt,
			sleep:        10 * time.Second,
			data:         &data,
			failed:       &failed,
		})
		if err != nil {
			return nil, err
		}

		// Check if instance is ready
		if ready, ok := data["ready"].(bool); ok && ready {
			data["id"] = instanceID
			return data, nil
		}

		// Not ready yet, sleep and retry
		tflog.Debug(ctx, fmt.Sprintf("Instance not ready yet, attempt=%d", attempt))
		attempt++
		select {
		case <-ctxTimeout.Done():
			return nil, fmt.Errorf("timeout reached while waiting for instance to be ready")
		case <-time.After(10 * time.Second):
			continue
		}
	}
}

func (api *API) waitUntilAllNodesReady(ctx context.Context, instanceID string) error {
	path := fmt.Sprintf("api/instances/%s/nodes", instanceID)
	ctxTimeout, cancel := context.WithTimeout(ctx, 1800*time.Second) // 30 minutes
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("waiting for all nodes to be ready, instanceID=%s", instanceID))
	attempt := 1

	for {
		if ctxTimeout.Err() != nil {
			return fmt.Errorf("timeout reached while waiting for all nodes to be ready")
		}

		var (
			data   []map[string]any
			failed map[string]any
		)

		tflog.Debug(ctx, fmt.Sprintf("Checking nodes ready status, attempt=%d", attempt))
		err := api.callWithRetry(ctxTimeout, api.sling.New().Get(path), retryRequest{
			functionName: "waitUntilAllNodesReady",
			resourceName: "Instance Nodes",
			attempt:      attempt,
			sleep:        15 * time.Second,
			data:         &data,
			failed:       &failed,
		})
		if err != nil {
			return err
		}

		tflog.Debug(ctx, fmt.Sprintf("response data=%v", data))

		// Check if all nodes are configured
		ready := true
		for _, node := range data {
			if configured, ok := node["configured"].(bool); ok {
				ready = ready && configured
			} else {
				ready = false
			}
		}

		if ready {
			return nil
		}

		// Not all nodes ready yet, sleep and retry
		tflog.Debug(ctx, fmt.Sprintf("Not all nodes ready yet, attempt=%d", attempt))
		attempt++
		select {
		case <-ctxTimeout.Done():
			return fmt.Errorf("timeout reached while waiting for all nodes to be ready")
		case <-time.After(15 * time.Second):
			continue
		}
	}
}

func (api *API) waitUntilAllNodesConfigured(ctx context.Context, instanceID string,
	attempt, sleep, timeout int) error {

	path := fmt.Sprintf("api/instances/%s/nodes", instanceID)
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("waiting for all nodes to be configured, instanceID=%s sleep=%d timeout=%d", instanceID, sleep, timeout))

	for {
		if ctxTimeout.Err() != nil {
			return fmt.Errorf("timeout reached after %d seconds, while waiting on all nodes configured", timeout)
		}

		var (
			data   []map[string]any
			failed map[string]any
		)

		tflog.Debug(ctx, fmt.Sprintf("Checking nodes configured status, attempt=%d", attempt))
		err := api.callWithRetry(ctxTimeout, api.sling.New().Get(path), retryRequest{
			functionName: "waitUntilAllNodesConfigured",
			resourceName: "Instance Nodes",
			attempt:      attempt,
			sleep:        time.Duration(sleep) * time.Second,
			data:         &data,
			failed:       &failed,
		})
		if err != nil {
			return err
		}

		tflog.Debug(ctx, fmt.Sprintf("response data=%v", data))

		// Check if all nodes are configured
		ready := true
		for _, node := range data {
			if configured, ok := node["configured"].(bool); ok {
				ready = ready && configured
			} else {
				ready = false
			}
		}

		if ready {
			return nil
		}

		// Not all nodes configured yet, sleep and retry
		tflog.Debug(ctx, fmt.Sprintf("Not all nodes configured yet, attempt=%d", attempt))
		attempt++
		select {
		case <-ctxTimeout.Done():
			return fmt.Errorf("timeout reached after %d seconds, while waiting on all nodes configured", timeout)
		case <-time.After(time.Duration(sleep) * time.Second):
			continue
		}
	}
}

func (api *API) waitUntilDeletion(ctx context.Context, instanceID string) error {
	path := fmt.Sprintf("/api/instances/%s", instanceID)
	ctxTimeout, cancel := context.WithTimeout(ctx, 1800*time.Second) // 30 minutes
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("waiting for instance deletion, instanceID=%s", instanceID))
	attempt := 1

	for {
		if ctxTimeout.Err() != nil {
			return fmt.Errorf("timeout reached while waiting for instance deletion")
		}

		var (
			data       map[string]any
			failed     map[string]any
			statusCode int
		)

		tflog.Debug(ctx, fmt.Sprintf("Checking instance deletion status, attempt=%d", attempt))
		err := api.callWithRetry(ctxTimeout, api.sling.New().Get(path), retryRequest{
			functionName: "waitUntilDeletion",
			resourceName: "Instance",
			attempt:      attempt,
			sleep:        10 * time.Second,
			data:         &data,
			failed:       &failed,
			statusCode:   &statusCode,
		})

		// Check if instance is deleted (404 or 410)
		if statusCode == 404 || statusCode == 410 {
			tflog.Debug(ctx, fmt.Sprintf("Instance deleted (status=%d)", statusCode))
			return nil
		}

		// If there was an error other than 404/410, return it
		if err != nil {
			return fmt.Errorf("failed to wait for deletion, error=%v", err)
		}

		// Instance still exists, sleep and retry
		tflog.Debug(ctx, fmt.Sprintf("Instance still exists, attempt=%d", attempt))
		attempt++
		select {
		case <-ctxTimeout.Done():
			return fmt.Errorf("timeout reached while waiting for instance deletion")
		case <-time.After(10 * time.Second):
			continue
		}
	}
}

func (api *API) CreateInstance(ctx context.Context, params map[string]any) (map[string]any, error) {
	var (
		data         map[string]any
		failed       map[string]any
		path         = "/api/instances"
		sensitiveCtx = tflog.MaskFieldValuesWithFieldKeys(ctx, "apikey", "url", "urls")
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s", path), params)
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateInstance",
		resourceName: "Instance",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	tflog.Debug(sensitiveCtx, "response data", data)
	if id, ok := data["id"]; ok {
		data["id"] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
	} else {
		return nil, fmt.Errorf("invalid identifier=%v", data["id"])
	}
	return api.waitUntilReady(ctx, data["id"].(string))
}

func (api *API) ReadInstance(ctx context.Context, instanceID string) (map[string]any, error) {
	var (
		data         map[string]any
		failed       map[string]any
		path         = fmt.Sprintf("/api/instances/%s", instanceID)
		sensitiveCtx = tflog.MaskFieldValuesWithFieldKeys(ctx, "apikey", "url", "urls")
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Path(path), retryRequest{
		functionName: "ReadInstance",
		resourceName: "Instance",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	// Handle resource drift
	if len(data) == 0 {
		return nil, nil
	}

	tflog.Debug(sensitiveCtx, "response data", data)
	return data, nil
}

func (api *API) UpdateInstance(ctx context.Context, instanceID string, params map[string]any) error {
	var (
		failed     map[string]any
		statusCode int
		path       = fmt.Sprintf("api/instances/%v", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s", path), params)
	err := api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateInstance",
		resourceName: "Instance",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
		statusCode:   &statusCode,
	})
	if err != nil {
		return err
	}

	// If resource was already deleted (410), skip waiting for nodes
	if statusCode == 410 {
		tflog.Debug(ctx, fmt.Sprintf("Instance already deleted (status=%d), skipping node readiness check", statusCode))
		return nil
	}

	return api.waitUntilAllNodesReady(ctx, instanceID)
}

func (api *API) DeleteInstance(ctx context.Context, instanceID string, keep_vpc bool) error {
	var (
		failed     map[string]any
		statusCode int
		path       = fmt.Sprintf("api/instances/%s?keep_vpc=%t", instanceID, keep_vpc)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteInstance",
		resourceName: "Instance",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
		statusCode:   &statusCode,
	})
	if err != nil {
		return err
	}

	// If resource was already deleted (404 or 410), skip waiting for deletion
	if statusCode == 404 || statusCode == 410 {
		tflog.Debug(ctx, fmt.Sprintf("Instance already deleted (status=%d), skipping deletion wait", statusCode))
		return nil
	}

	return api.waitUntilDeletion(ctx, instanceID)
}

func (api *API) UrlInformation(url string) map[string]any {
	paramsMap := make(map[string]any)
	r := regexp.MustCompile(`^.*:\/\/(?P<username>(.*)):(?P<password>(.*))@(?P<host>(.*))\/(?P<vhost>(.*))`)
	match := r.FindStringSubmatch(url)

	for i, value := range r.SubexpNames() {
		if value == "username" {
			paramsMap["username"] = match[i]
		}
		if value == "password" {
			paramsMap["password"] = match[i]
		}
		if value == "host" {
			paramsMap["host"] = match[i]
		}
		if value == "vhost" {
			paramsMap["vhost"] = match[i]
		}
	}

	return paramsMap
}

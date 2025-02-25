package provider

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// KarporClient is the Karpor client.
type KarporClient struct {
	Client      *http.Client
	ApiEndpoint string
	ApiKey      string
}

// NewKarporClient creates a new Karpor client.
func NewKarporClient(endpoint string, key string, skipTlsVerify bool) (*KarporClient, error) {
	return &KarporClient{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: skipTlsVerify,
				},
			},
			Timeout: 10 * time.Second,
		},
		ApiEndpoint: endpoint,
		ApiKey:      key,
	}, nil
}

// ValidateClusterConfig validates the cluster config.
func (c *KarporClient) ValidateClusterConfig(ctx context.Context, cluster *ClusterRegistrationResourceModel) (bool, error) {
	payloadData := map[string]string{
		"kubeConfig": cluster.Credentials.ValueString(),
	}
	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return false, err
	}
	payload := strings.NewReader(string(payloadBytes))

	req, err := http.NewRequest("POST", c.ApiEndpoint+"/rest-api/v1/cluster/config/validate", payload)
	if err != nil {
		return false, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return false, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return false, err
	}
	success, _ := data["success"].(bool)
	return success, nil
}

// RegisterCluster registers a new cluster.
func (c *KarporClient) RegisterCluster(ctx context.Context, cluster *ClusterRegistrationResourceModel) (string, error) {
	if cluster.DisplayName.IsNull() {
		cluster.DisplayName = types.StringValue(cluster.ClusterName.ValueString())
	}
	payloadData := map[string]string{
		"displayName": cluster.DisplayName.ValueString(),
		"description": cluster.Description.ValueString(),
		"kubeConfig":  cluster.Credentials.ValueString(),
	}
	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return "", err
	}
	payload := strings.NewReader(string(payloadBytes))

	req, err := http.NewRequest("POST", c.ApiEndpoint+"/rest-api/v1/cluster/"+cluster.ClusterName.ValueString(), payload)
	if err != nil {
		return "", err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}
	success, _ := data["success"].(bool)
	if !success {
		message, _ := data["message"].(string)
		return "", fmt.Errorf("%s", message)
	}
	clusterData, ok := data["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("missing data field in response")
	}

	metadata, ok := clusterData["metadata"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("missing metadata field in response")
	}

	uid, ok := metadata["uid"].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid uid field in response")
	}
	return uid, nil
}

// GetCluster gets a cluster.
func (c *KarporClient) GetCluster(ctx context.Context, clusterName string) (*ClusterRegistrationResourceModel, error) {
	req, err := http.NewRequest("GET", c.ApiEndpoint+"/rest-api/v1/cluster/"+clusterName, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	success, _ := data["success"].(bool)
	if !success {
		message, _ := data["message"].(string)
		return nil, fmt.Errorf("%s", message)
	}
	clusterData, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing data field in response")
	}
	metadata, ok := clusterData["metadata"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing metadata field in response")
	}

	uid, ok := metadata["uid"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid uid field in response")
	}
	returnedClusterName, ok := metadata["name"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid name field in response")
	}

	spec, ok := clusterData["spec"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing spec field in response")
	}

	displayName, ok := spec["displayName"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid displayName field in response")
	}

	description, ok := spec["description"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid description field in response")
	}

	remoteCluster := ClusterRegistrationResourceModel{
		Id:          types.StringValue(uid),
		ClusterName: types.StringValue(returnedClusterName),
		DisplayName: types.StringValue(displayName),
		Description: types.StringValue(description),
	}
	return &remoteCluster, nil
}

// UpdateCluster updates a cluster.
func (c *KarporClient) UpdateCluster(ctx context.Context, cluster *ClusterRegistrationResourceModel) (bool, error) {
	payloadData := map[string]string{
		"displayName": cluster.DisplayName.ValueString(),
		"description": cluster.Description.ValueString(),
	}
	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return false, err
	}
	payload := strings.NewReader(string(payloadBytes))

	req, err := http.NewRequest("PUT", c.ApiEndpoint+"/rest-api/v1/cluster/"+cluster.ClusterName.ValueString(), payload)
	if err != nil {
		return false, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return false, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return false, err
	}
	success, _ := data["success"].(bool)
	if !success {
		message, _ := data["message"].(string)
		return false, fmt.Errorf("%s", message)
	}
	return true, nil
}

// DeleteCluster deletes a cluster.
func (c *KarporClient) DeleteCluster(ctx context.Context, cluster *ClusterRegistrationResourceModel) (bool, error) {
	req, err := http.NewRequest("DELETE", c.ApiEndpoint+"/rest-api/v1/cluster/"+cluster.ClusterName.ValueString(), nil)
	if err != nil {
		return false, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return false, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return false, err
	}
	success, _ := data["success"].(bool)
	if !success {
		message, _ := data["message"].(string)
		return false, fmt.Errorf("%s", message)
	}
	return true, nil
}

func (c *KarporClient) doRequest(req *http.Request) ([]byte, error) {
	token := c.ApiKey

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

package aci

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/virtualcloudfoundry/goaci/api"
)

// ListContainerGroupUsage Get the usage for a subscription.
// From: https://docs.microsoft.com/en-us/rest/api/container-instances/containergroupusage/list
func (c *Client) ListContainerGroupUsage(location string) (*UsageListResult, error, *int) {
	urlParams := url.Values{
		"api-version": []string{apiVersion},
	}

	// Create the url.
	uri := api.ResolveRelative(BaseURI, containerGroupUsageURLPath)
	uri += "?" + url.Values(urlParams).Encode()

	// Create the request.
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("Creating list container group usage uri request failed: %v", err), nil
	}

	// Add the parameters to the url.
	if err := api.ExpandURL(req.URL, map[string]string{
		"subscriptionId": c.auth.SubscriptionID,
		"location":       location,
	}); err != nil {
		return nil, fmt.Errorf("Expanding URL with parameters failed: %v", err), nil
	}

	// Send the request.
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Sending list container group usage request failed: %v", err), &resp.StatusCode
	}
	defer resp.Body.Close()

	// 200 (OK) is a success response.
	if err := api.CheckResponse(resp); err != nil {
		return nil, err, &resp.StatusCode
	}

	// Decode the body from the response.
	if resp.Body == nil {
		return nil, errors.New("List container group usage returned an empty body in the response"), &resp.StatusCode
	}
	var cg UsageListResult
	if err := json.NewDecoder(resp.Body).Decode(&cg); err != nil {
		return nil, fmt.Errorf("Decoding list container group usage response body failed: %v", err), &resp.StatusCode
	}

	return &cg, nil, &resp.StatusCode
}

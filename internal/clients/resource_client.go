package clients

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const (
	moduleName    = "resource"
	moduleVersion = "v0.1.0"
)

type ResourceClient struct {
	host string
	pl   runtime.Pipeline
}

func NewResourceClient(credential azcore.TokenCredential, opt *arm.ClientOptions) (*ResourceClient, error) {
	if opt == nil {
		opt = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := opt.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, opt)
	if err != nil {
		return nil, err
	}
	return &ResourceClient{
		host: ep,
		pl:   pl,
	}, nil
}

func (client *ResourceClient) CreateOrUpdate(ctx context.Context, resourceID string, apiVersion string, body interface{}) (interface{}, error) {
	resp, err := client.createOrUpdate(ctx, resourceID, apiVersion, body)
	if err != nil {
		return nil, err
	}
	var responseBody interface{}
	pt, err := runtime.NewPoller[interface{}](resp, client.pl, nil)
	if err == nil {
		resp, err := pt.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 10 * time.Second,
		})
		return resp, err
	}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (client *ResourceClient) createOrUpdate(ctx context.Context, resourceID string, apiVersion string, body interface{}) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceID, apiVersion, body)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

func (client *ResourceClient) createOrUpdateCreateRequest(ctx context.Context, resourceID string, apiVersion string, body interface{}) (*policy.Request, error) {
	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, body)
}

func (client *ResourceClient) Get(ctx context.Context, resourceID string, apiVersion string) (interface{}, error) {
	req, err := client.getCreateRequest(ctx, resourceID, apiVersion)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return nil, runtime.NewResponseError(resp)
	}

	var responseBody interface{}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (client *ResourceClient) getCreateRequest(ctx context.Context, resourceID string, apiVersion string) (*policy.Request, error) {
	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

func (client *ResourceClient) Delete(ctx context.Context, resourceID string, apiVersion string) (interface{}, error) {
	resp, err := client.delete(ctx, resourceID, apiVersion)
	if err != nil {
		return nil, err
	}
	var responseBody interface{}
	pt, err := runtime.NewPoller[interface{}](resp, client.pl, nil)
	if err == nil {
		resp, err := pt.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 10 * time.Second,
		})
		return resp, err
	}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (client *ResourceClient) delete(ctx context.Context, resourceID string, apiVersion string) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceID, apiVersion)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

func (client *ResourceClient) deleteCreateRequest(ctx context.Context, resourceID string, apiVersion string) (*policy.Request, error) {
	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

func (client *ResourceClient) Action(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}) (interface{}, error) {
	resp, err := client.action(ctx, resourceID, action, apiVersion, method, body)
	if err != nil {
		return nil, err
	}
	var responseBody interface{}
	pt, err := runtime.NewPoller[interface{}](resp, client.pl, nil)
	if err == nil {
		resp, err := pt.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 10 * time.Second,
		})
		return resp, err
	}

	contentType := resp.Header.Get("Content-Type")
	switch {
	case strings.Contains(contentType, "text/plain"):
		payload, err := runtime.Payload(resp)
		if err != nil {
			return nil, err
		}
		responseBody = string(payload)
	case strings.Contains(contentType, "application/json"):
		if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
			return nil, err
		}
	default:
	}
	return responseBody, nil
}

func (client *ResourceClient) action(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}) (*http.Response, error) {
	req, err := client.actionCreateRequest(ctx, resourceID, action, apiVersion, method, body)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

func (client *ResourceClient) actionCreateRequest(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}) (*policy.Request, error) {
	urlPath := resourceID
	if action != "" {
		urlPath = fmt.Sprintf("%s/%s", resourceID, action)
	}
	req, err := runtime.NewRequest(ctx, method, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	if method != "GET" && body != nil {
		err = runtime.MarshalAsJSON(req, body)
	}
	return req, err
}

func (client *ResourceClient) List(ctx context.Context, url string, apiVersion string) (interface{}, error) {
	pager := runtime.NewPager[interface{}](runtime.PagingHandler[interface{}]{
		More: func(current interface{}) bool {
			if current == nil {
				return false
			}
			currentMap, ok := current.(map[string]interface{})
			if !ok {
				return false
			}
			if currentMap["nextLink"] == nil {
				return false
			}
			if nextLink := currentMap["nextLink"].(string); nextLink == "" {
				return false
			}
			return true
		},
		Fetcher: func(ctx context.Context, current *interface{}) (interface{}, error) {
			var request *policy.Request
			if current == nil {
				req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, url))
				if err != nil {
					return nil, err
				}
				reqQP := req.Raw().URL.Query()
				reqQP.Set("api-version", apiVersion)
				req.Raw().URL.RawQuery = reqQP.Encode()
				request = req
			} else {
				nextLink := ""
				if currentMap, ok := (*current).(map[string]interface{}); ok && currentMap["nextLink"] != nil {
					nextLink = currentMap["nextLink"].(string)
				}
				req, err := runtime.NewRequest(ctx, http.MethodGet, nextLink)
				if err != nil {
					return nil, err
				}
				request = req
			}
			request.Raw().Header.Set("Accept", "application/json")
			resp, err := client.pl.Do(request)
			if err != nil {
				return nil, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return nil, runtime.NewResponseError(resp)
			}
			var responseBody interface{}
			if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
				return nil, err
			}
			return responseBody, nil
		},
	})

	value := make([]interface{}, 0)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if pageMap, ok := page.(map[string]interface{}); ok {
			if pageMap["value"] != nil {
				if pageValue, ok := pageMap["value"].([]interface{}); ok {
					value = append(value, pageValue...)
					continue
				}
			}
		}

		// if response doesn't follow the ARM paging guideline, return the response as is
		return page, nil
	}
	return map[string]interface{}{
		"value": value,
	}, nil
}

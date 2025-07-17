package project

import (
	"context"
	"encoding/json"
	"fmt"
	model "github.com/dooray-go/dooray/openapi/model/project"
	"io"
	"net/http"
)

func (c *Project) GetPosts(apikey string, projectId string, toMemberIds string, postWorkflowClasses string) (*model.GetPostsResponse, error) {
	return c.GetPostsCustomHTTPContext(context.Background(), apikey, http.DefaultClient, projectId, toMemberIds, postWorkflowClasses)
}
func (c *Project) GetPostsContext(ctx context.Context, apikey string, projectId string, toMemberIds string, postWorkflowClasses string) (*model.GetPostsResponse, error) {
	return c.GetPostsCustomHTTPContext(ctx, apikey, http.DefaultClient, projectId, toMemberIds, postWorkflowClasses)
}
func (c *Project) GetPostsCustomHTTP(apikey string, httpClient *http.Client, projectId string, toMemberIds string, postWorkflowClasses string) (*model.GetPostsResponse, error) {
	return c.GetPostsCustomHTTPContext(context.Background(), apikey, httpClient, projectId, toMemberIds, postWorkflowClasses)
}

func (c *Project) GetPostsCustomHTTPContext(ctx context.Context, apikey string, httpClient *http.Client, projectId string, toMemberIds string, postWorkflowClasses string) (*model.GetPostsResponse, error) {
	url := c.endPoint + "/project/v1/projects/" + projectId + "/posts"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()

	if toMemberIds != "" {
		query.Add("toMemberIds", toMemberIds)
	}

	if postWorkflowClasses != "" {
		query.Add("postWorkflowClasses", postWorkflowClasses)
	}

	req.URL.RawQuery = query.Encode()

	req.Header.Set("Authorization", "dooray-api "+apikey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		resp.Body.Close()
	}()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var getpostsResponse model.GetPostsResponse
	if err := json.Unmarshal(resBody, &getpostsResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	getpostsResponse.RawJSON = string(resBody)

	return &getpostsResponse, nil
}

package project

import (
	"context"
	"io"
	"net/http"
)

func (c *Project) GetPosts(apikey string, projectId string, toMemberIds string) (string, error) {
	return c.GetPostsCustomHTTPContext(context.Background(), apikey, http.DefaultClient, projectId, toMemberIds)
}
func (c *Project) GetPostsContext(ctx context.Context, apikey string, projectId string, toMemberIds string) (string, error) {
	return c.GetPostsCustomHTTPContext(ctx, apikey, http.DefaultClient, projectId, toMemberIds)
}
func (c *Project) GetPostsCustomHTTP(apikey string, httpClient *http.Client, projectId string, toMemberIds string) (string, error) {
	return c.GetPostsCustomHTTPContext(context.Background(), apikey, httpClient, projectId, toMemberIds)
}

func (c *Project) GetPostsCustomHTTPContext(ctx context.Context, apikey string, httpClient *http.Client, projectId string, toMemberIds string) (string, error) {
	url := c.endPoint + "/project/v1/projects/" + projectId + "/posts"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	query := req.URL.Query()

	if toMemberIds != "" {
		query.Add("toMemberIds", toMemberIds)
	}

	req.URL.RawQuery = query.Encode()

	req.Header.Set("Authorization", "dooray-api "+apikey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		resp.Body.Close()
	}()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(resBody), nil
}

package account

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (a *Account) GetMember(apikey string, id string) (string, error) {
	return a.GetMemberCustomHTTPContext(context.Background(), apikey, http.DefaultClient, id)
}

func (a *Account) GetMemberContext(ctx context.Context, apikey string, id string) (string, error) {
	return a.GetMemberCustomHTTPContext(ctx, apikey, http.DefaultClient, id)
}

func (a *Account) GetMemberCustomHTTP(apikey string, httpClient *http.Client, id string) (string, error) {
	return a.GetMemberCustomHTTPContext(context.Background(), apikey, httpClient, id)
}

func (a *Account) GetMemberCustomHTTPContext(ctx context.Context, apikey string, httpClient *http.Client, id string) (string, error) {

	url := a.endPoint + "/common/v1/members/" + id
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed new request: %w", err)
	}

	query := req.URL.Query()
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Authorization", fmt.Sprintf("dooray-api %s", apikey))

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get member: %w", err)
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

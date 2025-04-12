package account

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Account struct {
	endPoint string
}

func NewDefaultAccount() *Account {
	return &Account{
		endPoint: "https://api.dooray.com",
	}
}

func NewAccount(endPoint string) *Account {
	return &Account{
		endPoint: endPoint,
	}
}

func (a *Account) GetMembers(apikey string, name string, userCode string) (string, error) {
	return a.GetMembersCustomHTTPContext(context.Background(), apikey, http.DefaultClient, name, userCode)
}

func (a *Account) GetMembersContext(ctx context.Context, apikey string, name string, userCode string) (string, error) {
	return a.GetMembersCustomHTTPContext(ctx, apikey, http.DefaultClient, name, userCode)
}

func (a *Account) GetMembersCustomHTTP(apikey string, httpClient *http.Client, name string, userCode string) (string, error) {
	return a.GetMembersCustomHTTPContext(context.Background(), apikey, httpClient, name, userCode)
}

func (a *Account) GetMembersCustomHTTPContext(ctx context.Context, apikey string, httpClient *http.Client, name string, userCode string) (string, error) {

	url := a.endPoint + "/common/v1/members"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed new request: %w", err)
	}

	query := req.URL.Query()

	if name != "" {
		query.Add("name", name)
	}

	if userCode != "" {
		query.Add("userCode", userCode)
	}

	req.URL.RawQuery = query.Encode()

	req.Header.Set("Authorization", fmt.Sprintf("dooray-api %s", apikey))

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get members: %w", err)
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

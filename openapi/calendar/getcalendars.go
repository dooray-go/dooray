package calendar

import (
	"context"
	"io"
	"net/http"
)

func (c *Calendar) GetCalendars(apikey string) (string, error) {
	return c.GetCalendarsCustomHTTPContext(context.Background(), apikey, http.DefaultClient)
}
func (c *Calendar) GetCalendarsContext(ctx context.Context, apikey string) (string, error) {
	return c.GetCalendarsCustomHTTPContext(ctx, apikey, http.DefaultClient)
}
func (c *Calendar) GetCalendarsCustomHTTP(apikey string, httpClient *http.Client) (string, error) {
	return c.GetCalendarsCustomHTTPContext(context.Background(), apikey, httpClient)
}
func (c *Calendar) GetCalendarsCustomHTTPContext(ctx context.Context, apikey string, httpClient *http.Client) (string, error) {
	url := c.endPoint + "/calendar/v1/calendars"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	query := req.URL.Query()
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

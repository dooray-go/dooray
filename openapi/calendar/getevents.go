package calendar

import (
	"context"
	"io"
	"net/http"
	"time"
)

func (c *Calendar) GetEvents(apikey string, calendars string, timeMin time.Time, timeMax time.Time) (string, error) {
	return c.GetEventsCustomHTTPContext(context.Background(), apikey, http.DefaultClient, calendars, timeMin, timeMax)
}
func (c *Calendar) GetEventsContext(ctx context.Context, apikey string, calendars string, timeMin time.Time, timeMax time.Time) (string, error) {
	return c.GetEventsCustomHTTPContext(ctx, apikey, http.DefaultClient, calendars, timeMin, timeMax)
}
func (c *Calendar) GetEventsCustomHTTP(apikey string, httpClient *http.Client, calendars string, timeMin time.Time, timeMax time.Time) (string, error) {
	return c.GetEventsCustomHTTPContext(context.Background(), apikey, httpClient, calendars, timeMin, timeMax)
}

func (c *Calendar) GetEventsCustomHTTPContext(ctx context.Context, apikey string, httpClient *http.Client, calendars string,
	timeMin time.Time, timeMax time.Time) (string, error) {
	url := c.endPoint + "/calendar/v1/calendars/*/events"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	query := req.URL.Query()

	if calendars != "" {
		query.Add("calendars", calendars)
	}

	query.Add("timeMin", FormatTimeToISO8601(timeMin))
	query.Add("timeMax", FormatTimeToISO8601(timeMax))

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

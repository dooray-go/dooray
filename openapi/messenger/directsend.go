package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dooray-go/dooray/utils"
	"io"
	"net/http"
)

type Messenger struct {
	endPoint string
}

func NewDefaultMessenger() *Messenger {
	return &Messenger{
		endPoint: "https://api.dooray.com",
	}
}

func NewMessenger(endPoint string) *Messenger {
	return &Messenger{
		endPoint: endPoint,
	}
}

type DirectSendRequest struct {
	Text                 string `json:"text"`
	OrganizationMemberId string `json:"organizationMemberId"`
}

func (m Messenger) DirectSend(apikey string, msg *DirectSendRequest) error {
	return m.DirectSendCustomHTTPContext(context.Background(), apikey, http.DefaultClient, msg)
}

func (m Messenger) DirectSendContext(ctx context.Context, apikey string, msg *DirectSendRequest) error {
	return m.DirectSendCustomHTTPContext(ctx, apikey, http.DefaultClient, msg)
}

func (m Messenger) DirectSendCustomHTTP(apikey string, httpClient *http.Client, msg *DirectSendRequest) error {
	return m.DirectSendCustomHTTPContext(context.Background(), apikey, httpClient, msg)
}

func (m Messenger) DirectSendCustomHTTPContext(ctx context.Context, apikey string, httpClient *http.Client, msg *DirectSendRequest) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	url := m.endPoint + "/messenger/v1/channels/direct-send"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("failed new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("dooray-api %s", apikey))

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post webhook: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	return checkStatusCode(resp)
}

func checkStatusCode(resp *http.Response) error {

	if resp.StatusCode != http.StatusOK {
		return utils.StatusCodeError{Code: resp.StatusCode, Status: resp.Status}
	}

	return nil
}

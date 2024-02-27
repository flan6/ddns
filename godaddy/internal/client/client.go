package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/flan6/ddns/godaddy/entity"
)

type GoDaddyClient struct {
	httpClient *http.Client
	secret     entity.Secret
}

// NewGodaddyClient facilitates dependency injection
// as each GodaddyClient property must be passed down to the function
func NewGodaddyClient(c *http.Client, s entity.Secret) GoDaddyClient {
	return GoDaddyClient{c, s}
}

func (c GoDaddyClient) Request(method, url string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("sso-key %s:%s", c.secret.ApiKey, c.secret.ApiSecret))

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		e := entity.GoDaddyError{}
		if err := json.NewDecoder(response.Body).Decode(&e); err != nil {
			return nil, err
		}

		e.StatusCode = response.StatusCode

		return nil, e
	}

	return io.ReadAll(response.Body)
}

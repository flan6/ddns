package godaddy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flan6/ddns/godaddy/entity"
	"github.com/flan6/ddns/godaddy/internal/client"
)

const (
	GoDaddyURL = "https://api.godaddy.com/v1/"
)

type Client interface {
	Request(method, path string, data []byte) ([]byte, error)
}

func NewGoDaddy(apiKey, apiSecret string) GoDaddy {
	client := client.NewGodaddyClient(
		http.DefaultClient,
		entity.Secret{
			ApiKey:    apiKey,
			ApiSecret: apiSecret,
		},
	)

	return GoDaddy{client}
}

type GoDaddy struct {
	c Client
}

func (g GoDaddy) Domains() ([]entity.Domain, error) {
	response, err := g.c.Request(http.MethodGet, fmt.Sprint(GoDaddyURL, "domains"), nil)
	if err != nil {
		return nil, err
	}

	result := make([]entity.Domain, 0)
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (g GoDaddy) Records(domainName string) ([]entity.Record, error) {
	response, err := g.c.Request(
		http.MethodGet,
		fmt.Sprintf("%sdomains/%s/records", GoDaddyURL, domainName),
		nil,
	)
	if err != nil {
		return nil, err
	}

	result := make([]entity.Record, 0)
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (g GoDaddy) RecordsByType(domainName, recordType, recordName string) ([]entity.Record, error) {
	response, err := g.c.Request(
		http.MethodGet,
		fmt.Sprintf("%sdomains/%s/records/%s/%s", GoDaddyURL, domainName, recordType, recordName),
		nil,
	)
	if err != nil {
		return nil, err
	}

	result := make([]entity.Record, 0)
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Replace all DNS Records for the specified Domain with the specified recordType and recordName
func (g GoDaddy) SetRecordsByType(domainName, recordType, recordName string, records []entity.Record) error {
	data, err := json.Marshal(records)
	if err != nil {
		return err
	}

	_, err = g.c.Request(
		http.MethodPut,
		fmt.Sprintf("%sdomains/%s/records/%s/%s", GoDaddyURL, domainName, recordType, recordName),
		data,
	)

	return err
}

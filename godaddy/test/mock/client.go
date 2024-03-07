package mock

import "errors"

// ClientMock is used for testing implementing the Client interface.
type ClientMock struct {
	RemoteURL string
}

func (c ClientMock) Request(method, path string, data []byte) ([]byte, error) {
	switch path {
	case c.RemoteURL + "domains":
		return []byte(`[{"domain": "example.com"}]`), nil

	case c.RemoteURL + "domains/test/records":
		return []byte(`[{"data": "255.255.255.255"}]`), nil

	default:
		return nil, errors.New("path not found")
	}
}

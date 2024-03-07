package godaddy

import (
	"testing"

	"github.com/flan6/ddns/godaddy/test/mock"
)

func TestDomains(t *testing.T) {
	clientMock := mock.ClientMock{RemoteURL: GoDaddyURL}

	daddy := GoDaddy{clientMock}

	_, err := daddy.Domains()
	if err != nil {
		t.Errorf("Domains() got err: %s", err)
	}
}

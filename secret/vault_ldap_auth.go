package secret

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/briandowns/spinner"
	httpRequests "github.com/rnkoaa/vault-env/http"
)

var s = spinner.New(spinner.CharSets[60], 100*time.Millisecond) // Build our new spinner

// VaultAuth -
type VaultAuth struct {
	Method   string
	Username string
	Password string
	Token    string
}

// NewAuth - create a new vaultLdapAuth object
func NewAuth(method, username, password, token string) *VaultAuth {
	return &VaultAuth{
		Method:   method,
		Username: username,
		Password: password,
		Token:    token,
	}
}

// VaultResponseAuth -
type VaultResponseAuth struct {
	ClientToken   string            `json:"client_token,omitempty"`
	Accessor      string            `json:"accessor,omitempty"`
	Policies      []string          `json:"policies,omitempty"`
	LeaseDuration int               `json:"lease_duration,omitempty"`
	Renewable     bool              `json:"renewable,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	EntityID      string            `json:"entity_id,omitempty"`
}

// VaultResponse -
type VaultResponse struct {
	RequestID     string             `json:"request_id,omitempty"`
	LeaseID       string             `json:"lease_id,omitempty"`
	Renewable     bool               `json:"renewable,omitempty"`
	LeaseDuration int                `json:"lease_duration,omitempty"`
	Auth          *VaultResponseAuth `json:"auth,omitempty"`
}

// ToString -
func (v *VaultResponseAuth) ToString() string {
	return fmt.Sprintf("Token: %s, LeaseDuration: %d", v.ClientToken, v.LeaseDuration)
}

func processAuthResponse(resp *http.Response) (*VaultResponseAuth, error) {
	if s != nil {
		s.Stop()
	}
	if httpRequests.IsSuccessful(resp) {
		defer resp.Body.Close()
		rb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var response VaultResponse
		if rb != nil {
			err = json.Unmarshal(rb, &response)
			if err != nil {
				return nil, err
			}
		}
		return response.Auth, nil
	}
	// not successful
	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, fmt.Errorf("error not found. check request and try again")
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("error unauthorized, please retry the request with proper credentials")
	case http.StatusForbidden:
		return nil, fmt.Errorf("%d token expired", http.StatusForbidden)
	}
	return nil, nil
}

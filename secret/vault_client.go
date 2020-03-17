package secret

import (
	"fmt"
	"log"

	"github.com/hashicorp/vault/api"
	httpRequests "github.com/rnkoaa/vault-env/http"
)

// VaultClient - client to vault
type VaultClient struct {
	IsAuthenticated bool
	Address         string
	Client          *api.Client
	Auth            *VaultAuth
}

// NewClient -
func NewClient(address string, v *VaultAuth) *VaultClient {
	if v == nil {
		v = NewAuth(defaultAuthMethod, "", "", defaultToken)
	}

	if address == "" {
		address = defaultVaultAddress
	}

	return &VaultClient{
		Address: address,
		Client:  initVaultClient(address, v),
		Auth:    v,
	}
}

// Login -
func (v *VaultClient) Login() (*VaultResponseAuth, error) {
	auth := v.Auth
	log.Printf("\nlogging into enterprise secrets for %s", auth.Username)
	s.Start()
	requestURL := fmt.Sprintf("%s/v1/auth/ldap/login/%s", v.Address, auth.Username)
	body := fmt.Sprintf("{\"password\": \"%s\"}", auth.Password)

	reqBody := []byte(body)
	resp, err := httpRequests.MakePostForResponse(requestURL, reqBody)
	if err != nil {
		log.Printf("error retrieving user auth code, %#v", err.Error())
		return nil, err
	}

	return processAuthResponse(resp)
}

// RenewToken Renew a Token that has not expired
func (v *VaultClient) RenewToken(token string) (*VaultResponseAuth, error) {
	requestURL := fmt.Sprintf("%s/v1/auth/token/renew-self", v.Address)
	body := `{"increment": 1800}`
	reqBody := []byte(body)
	reqHeaders := make(map[string]interface{})
	reqHeaders["X-Vault-Token"] = token
	resp, err := httpRequests.MakePostWithHeadersForResponse(requestURL, reqHeaders, reqBody)
	if err != nil {
		return nil, err
	}
	vResponse, err := processAuthResponse(resp)
	if err != nil {
		return nil, err
	}

	if vResponse != nil {
		v.Auth.Token = vResponse.ClientToken
	}
	return vResponse, nil
}

func initVaultClient(address string, v *VaultAuth) *api.Client {
	conf := &api.Config{
		Address: address,
	}
	c, err := api.NewClient(conf)
	if err != nil {
		log.Fatalf(err.Error())
	}
	c.SetToken(v.Token)
	return c
}

package secret

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/vault/api"
	httpRequests "github.com/rnkoaa/vault-env/http"
)

// VaultClient - client to vault
type VaultClient struct {
	IsAuthenticated bool
	SecretsVersion  int
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
func (v *VaultClient) Login() error {
	switch v.Auth.Method {
	case strings.ToLower("Token"):
		return processTokenAuth(v)
	case strings.ToLower("Ldap"):
		return processLDAPLogin(v)
	}

	return nil
}

func processTokenAuth(v *VaultClient) error {
	if v.Auth.Token == "" {
		return errorEmptyToken
	}
	v.IsAuthenticated = true
	v.Client.SetToken(v.Auth.Token)
	return nil
}

func processLDAPLogin(v *VaultClient) error {
	authResponse, err := ldapLogin(v.Address, v.Auth)
	if err != nil {
		return err
	}
	if authResponse == nil {
		return fmt.Errorf("no vault response")
	}
	v.Auth.Token = authResponse.ClientToken
	v.IsAuthenticated = true
	v.Client.SetToken(authResponse.ClientToken)
	return nil
}

func ldapLogin(address string, auth *VaultAuth) (*VaultResponseAuth, error) {
	log.Printf("logging into vault using ldap for user %s", auth.Username)
	s.Start()
	requestURL := fmt.Sprintf("%s/v1/auth/ldap/login/%s", address, auth.Username)
	body := fmt.Sprintf(`{"password": "%s"}`, auth.Password)

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
	reqBody := []byte(`{"increment": 1800}`)
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

// ResolveValues -
func (v *VaultClient) ResolveValues(values map[string]interface{}) (map[string]interface{}, map[string]error) {
	errs := make(map[string]error, 0)
	res := make(map[string]interface{})
	if !v.IsAuthenticated {
		err := v.Login()
		if err != nil {
			errs["AuthError"] = err
			return res, errs
		}
	}
	for k, value := range values {
		if v != nil {
			value, err := v.ResolveValue(value.(string))
			if err != nil {
				errs[k] = err
				if len(errs) >= 3 {
					log.Printf("too many errors encountered when reading secrets. exiting...")
					break
				}
			} else {
				res[k] = value
			}
		} else {
			errs[k] = errors.New("nil value for key")
			if len(errs) >= 3 {
				log.Printf("too many errors encountered when reading secrets. exiting...")
				break
			}
		}
	}
	return res, errs
}

func processKey(secretsVersion int, key string) string {
	if strings.HasPrefix(key, secretPathPrefix) {
		return key
	}
	switch secretsVersion {
	case 1:
		return fmt.Sprintf(v1SecretPathPrefix, key)
	default:
		return fmt.Sprintf(v2SecretPathPrefix, key)
	}
}

// ResolveValue -
func (v *VaultClient) ResolveValue(key string) (string, error) {
	if !v.IsAuthenticated {
		err := v.Login()
		if err != nil {
			return "", err
		}
	}

	// process the key path based on the version of secrets engine.
	k := processKey(v.SecretsVersion, key)
	secret, err := v.Client.Logical().Read(k)
	if err != nil {
		return "", err
	}
	if secret == nil {
		return "", fmt.Errorf("nil secret retrieved")
	}
	if secret.Data != nil && len(secret.Data) > 0 {
		switch v.SecretsVersion {
		case 1:
			return processV1SecretDataResponse(secret)
		default:
			return processV2SecretDataResponse(secret)
		}
	}

	return "", nil
}

func processV2SecretDataResponse(secret *api.Secret) (string, error) {
	if data, dataOk := secret.Data["data"]; dataOk {
		secretData := data.(map[string]interface{})
		if secretDataValue, ok := secretData["value"]; ok && secretDataValue != nil {
			return secretDataValue.(string), nil
		}
	}
	return "", nil
}

func processV1SecretDataResponse(secret *api.Secret) (string, error) {
	if secretDataValue, ok := secret.Data["value"]; ok && secretDataValue != nil {
		return secretDataValue.(string), nil
	}
	return "", nil
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

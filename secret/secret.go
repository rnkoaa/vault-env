package secret

import "fmt"

var (
	defaultVaultAddress = "localhost:8200"
	defaultAuthMethod   = "Token"
	defaultToken        = "default-token"
	secretPathPrefix    = "secret/"
	v1SecretPathPrefix  = "secret/%s"
	v2SecretPathPrefix  = "secret/data/%s"
	errorEmptyToken     = fmt.Errorf("token method desired but token is empty")
)

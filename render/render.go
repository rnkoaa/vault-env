package render

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/rnkoaa/vault-env/dirs"
	"github.com/rnkoaa/vault-env/helpers"
	"github.com/rnkoaa/vault-env/secret"
)

var (
	// ErrorInputFileNotFound -
	ErrorInputFileNotFound  = fmt.Errorf("Input file not found")
	errorTooManyVaultErrors = errors.New("too many vault errors")
)

// Config - config to process
type Config struct {
	InputFile    string
	InFormat     string
	OutputFile   string
	OutputFormat string
}

// VaultConfig -
type VaultConfig struct {
	URL           string
	SecretVersion int
	AuthEnabled   bool
	AuthMethod    string
	Username      string
	Password      string
	Token         string
	AuthFile      string // file containing username and password to authenticate to vault as well as url
}

// New - create a new render config
func New(inFormat, inFile, outFile, outFormat string) *Config {
	return &Config{
		InFormat:     inFormat,
		InputFile:    inFile,
		OutputFile:   outFile,
		OutputFormat: outFormat,
	}
}

// ToString -
func (r *Config) ToString() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

// ToString -
func (c *VaultConfig) ToString() string {
	res := ""
	if c != nil {
		if c.URL != "" {
			res += fmt.Sprintf("Url: %s", c.URL)
		}
		if c.AuthEnabled {
			if c.AuthMethod != "" {
				if res == "" {
					res += fmt.Sprintf("AuthMethod: %s", c.AuthMethod)
				} else {
					res += fmt.Sprintf(",\nAuthMethod: %s", c.AuthMethod)
				}
			}
			if c.Token != "" {
				if res == "" {
					res += fmt.Sprintf("Token: %s", c.Token)
				} else {
					res += fmt.Sprintf(",\nToken: %s", c.Token)
				}
			}
			if c.Password != "" {
				if res == "" {
					res += fmt.Sprintf("Password: %s", c.Password)
				} else {
					res += fmt.Sprintf(",\nPassword: %s", c.Password)
				}
			}
			if c.Username != "" {
				if res == "" {
					res += fmt.Sprintf("Username: %s", c.Username)
				} else {
					res += fmt.Sprintf(",\nUsername: %s", c.Username)
				}
			}
			if res == "" {
				res += fmt.Sprintf("AuthEnabled: %v", c.AuthEnabled)
			} else {
				res += fmt.Sprintf(",\nAuthEnabled: %v", c.AuthEnabled)
			}
		}
	}
	return res
}

// CreateVaultClient -
func CreateVaultClient(v *VaultConfig) *secret.VaultClient {
	auth := createAuth(v)
	client := secret.NewClient(v.URL, auth)
	client.SecretsVersion = v.SecretVersion
	return client
}

func createAuth(v *VaultConfig) *secret.VaultAuth {
	return &secret.VaultAuth{
		Method:   v.AuthMethod,
		Username: v.Username,
		Token:    v.Token,
		Password: v.Password,
	}
}

// DoRender - the actual rendering of the secret file happens here.
func DoRender(v *secret.VaultClient, r *Config) error {
	if r.InputFile == "" {
		return ErrorInputFileNotFound
	}
	if r.InFormat == "" {
		inputFormat, err := helpers.GetFileType(r.InputFile)
		if err != nil {
			return err
		}
		r.InFormat = inputFormat
	}

	// file exists and can be determined.
	switch r.InFormat {
	case "yaml":
		return processYamlFile(v, r)
	case "json":
		return processJSONFile(v, r)
	}
	return nil
}

func processYamlFile(v *secret.VaultClient, c *Config) error {
	d, err := ioutil.ReadFile(c.InputFile)
	if err != nil {
		return err
	}
	flattened := helpers.FlattenYaml(d)
	if flattened == nil {
		return errors.New("error flattening yaml, look into fixing it")
	}

	values, errs := v.ResolveValues(flattened)
	if len(errs) > 0 {
		for k, e := range errs {
			log.Printf("error resolving key: %s, %v\n", k, e)
		}

		return errorTooManyVaultErrors
	}
	resolvedValues := make(map[string]interface{}, len(flattened))
	for k, v := range flattened {
		if v != nil {
			if secret, ok := values[k]; ok {
				resolvedValues[k] = secret
			}
		}
	}

	// expand resolvedValues
	res, err := helpers.ExpandYaml(resolvedValues)
	if err != nil {
		return err
	}
	// write to file
	// write the content overwriting any content of the file.
	return dirs.WriteContentToFile(c.OutputFile, res, true)
}

func processJSONFile(v *secret.VaultClient, c *Config) error {
	return nil
}

func printMap(m map[string]interface{}) {
	for k, v := range m {
		if v != nil {
			fmt.Printf("Key: %s, Value: %s\n", k, v.(string))
		} else {
			fmt.Printf("Key: %s has null value.\n", k)
		}
	}
}

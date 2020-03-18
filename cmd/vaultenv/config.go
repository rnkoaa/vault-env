package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// VaultConfig -
type VaultConfig struct {
	URL           string
	SecretVersion int
	AuthEnabled   bool
	AuthMethod    string
	Username      string
	Password      string
	Token         string
}

// LoadConfig - create a new config object from viper configs
// mostly from the config.yml and/or command line args.
func LoadConfig() *VaultConfig {
	conf := &VaultConfig{}
	conf.URL = viper.GetString("vault.url")
	conf.AuthEnabled = viper.GetBool("vault.auth.enabled")
	conf.Username = viper.GetString("vault.auth.username")
	conf.AuthMethod = viper.GetString("vault.auth.method")
	conf.Password = viper.GetString("vault.auth.password")
	conf.Token = viper.GetString("vault.auth.token")
	secretsVersion := viper.GetInt("secrets.version")

	// secretsVersion paths changed from version 1 to version 2
	if secretsVersion < 1 {
		secretsVersion = 1
	}

	conf.SecretVersion = secretsVersion
	return conf
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

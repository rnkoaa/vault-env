package main

import (
	"github.com/rnkoaa/vault-env/render"
	"github.com/spf13/viper"
)

// LoadConfig - create a new config object from viper configs
// mostly from the config.yml and/or command line args.
func LoadConfig() *render.VaultConfig {
	conf := &render.VaultConfig{}
	conf.URL = viper.GetString("vault.url")
	conf.AuthEnabled = viper.GetBool("vault.auth.enabled")
	conf.Username = viper.GetString("vault.auth.username")
	conf.AuthMethod = viper.GetString("vault.auth.method")
	conf.Password = viper.GetString("vault.auth.password")
	conf.Token = viper.GetString("vault.auth.token")
	secretsVersion := viper.GetInt("vault.secrets.version")

	// secretsVersion paths changed from version 1 to version 2
	if secretsVersion < 1 {
		secretsVersion = 1
	}

	conf.SecretVersion = secretsVersion
	return conf
}

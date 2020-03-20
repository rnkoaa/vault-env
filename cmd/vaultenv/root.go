package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Execute -
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is config.yml)")

	// renderCmd.PersistentFlags().StringVarP(&secretTemplateFile, "input", "i", "vault.yml", "input file to process for vault keys")
	renderCmd.PersistentFlags().StringP("input", "i", "vault.yml", "input file to process for vault keys")
	renderCmd.PersistentFlags().StringP("output", "o", "secret.yml", "output file of output (yaml, json, table)")
	renderCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", defaultFormat, "format of output (yaml, json, table)")
	viper.BindPFlag("secrets.output.file", renderCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("secret.template.file", renderCmd.PersistentFlags().Lookup("input"))
	rootCmd.AddCommand(renderCmd)

	// outputFile = viper.GetString("secrets.output.file")
	initConfig()
}

func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}
	if err := viper.ReadInConfig(); err != nil {
		er("error reading config file. please provide config.")
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile("runtime.yml")
	if err := viper.MergeInConfig(); err != nil {
		fmt.Println("error reading and merging config file. please provide config.")
		er(err)
	}

	vaultAuthFile := viper.GetString("vault.auth.file")
	if vaultAuthFile != "" {
		viper.SetConfigType("yaml")
		viper.SetConfigFile(vaultAuthFile)
		if err := viper.MergeInConfig(); err != nil {
			fmt.Println("error reading and merging vault auth file. ensure it exists")
			er(err)
		}
	}

	// viper.AddConfigPath("/tap/secret/secret")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	vaultConf = LoadConfig()
	viperKeys := viper.AllKeys()
	for _, key := range viperKeys {
		fmt.Printf("Key: %s\n", key)
	}
}

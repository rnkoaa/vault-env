package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	// https://github.com/spf13/cobra
	// https://github.com/spf13/viper
	configFile   string
	outputFile   string
	inputFile    string
	outputFormat string
	userLicense  string

	defaultFormat = "yaml"
	vaultConf     *VaultConfig
	rootCmd       = &cobra.Command{
		Use:   "cobra",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
)

func main() {

	Execute()
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is config.yml)")

	renderCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "vault.yml", "input file to process for vault keys")
	renderCmd.PersistentFlags().StringP("output", "o", "secret.yml", "output file of output (yaml, json, table)")
	renderCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", defaultFormat, "format of output (yaml, json, table)")
	viper.BindPFlag("secrets.output.file", renderCmd.PersistentFlags().Lookup("output"))
	rootCmd.AddCommand(renderCmd)

	outputFile = viper.GetString("secrets.output.file")
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

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		er("error reading config file. please provide config.")
	}
	vaultConf = LoadConfig()

}

package main

import (
	"fmt"
	"os"

	"github.com/rnkoaa/vault-env/render"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	configFile         string
	outputFile         string
	secretTemplateFile string
	outputFormat       string
	userLicense        string

	defaultFormat = "yaml"
	vaultConf     *render.VaultConfig
	rootCmd       = &cobra.Command{
		Use:   "cobra",
		Short: "A generator for Cobra based Applications",
		Long:  `Cobra is a CLI library for Go that empowers applications. This application is a tool to generate the needed files to quickly create a Cobra application.`,
	}
)

func main() {

	Execute()
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

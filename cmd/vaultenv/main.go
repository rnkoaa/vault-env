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
		Use:   "vault-env",
		Short: "An application to sync enterprise secrets vault into containers before they start up.",
		Long: `As applications are deployed to in containers, they will need their secrets. 
In order to prevent secret sprawl and prevent secret leakage, 
this application takes authentication configs from the orchestration environment, 
logs into vault and generates the secret properties before the application starts up.`,
	}
)

func main() {

	Execute()
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

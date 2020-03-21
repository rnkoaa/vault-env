package main

import (
	"fmt"
	"os"

	"github.com/rnkoaa/vault-env/render"
	"github.com/rnkoaa/vault-env/secret"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	errorUnsupportedFileType = fmt.Errorf("unsupported file type")
	vaultClient              *secret.VaultClient
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "render values from vault",
	Long:  `render values from vault using a configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		secretTemplateFile = viper.GetString("vault.secret.template.file")
		outputFile = viper.GetString("vault.secret.output.file")
		if secretTemplateFile == "" {
			fmt.Printf("secret template file '<empty>' not found\n")
			er(render.ErrorInputFileNotFound)
		}
		renderConfig := render.New(defaultFormat, secretTemplateFile, outputFile, outputFormat)
		if renderConfig == nil {
			fmt.Println("error creating renderConfig")
			os.Exit(1)
		}
		vaultClient = render.CreateVaultClient(vaultConf)

		if err := render.DoRender(vaultClient, renderConfig); err != nil {
			er(err)
		} else {
			fmt.Printf("Done rendering secrets to %s\n", outputFile)
		}
	},
}

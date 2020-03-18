package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/rnkoaa/vault-env/dirs"
	"github.com/rnkoaa/vault-env/helpers"
	"github.com/rnkoaa/vault-env/secret"
	"github.com/spf13/cobra"
)

var (
	errorInputFileNotFound   = fmt.Errorf("Input file not found")
	errorUnsupportedFileType = fmt.Errorf("unsupported file type")
	errorTooManyVaultErrors  = errors.New("too many vault errors")
	vaultClient              *secret.VaultClient
)

// RenderConfig - config to process
type RenderConfig struct {
	InputFile    string
	InFormat     string
	OutputFile   string
	OutputFormat string
}

func new(inFormat, inFile, outFile, outFormat string) *RenderConfig {
	return &RenderConfig{
		InFormat:     inFormat,
		InputFile:    inFile,
		OutputFile:   outFile,
		OutputFormat: outFormat,
	}
}

// ToString -
func (r *RenderConfig) ToString() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "render values from vault",
	Long:  `render values from vault using a configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		vaultClient = createVaultClient(vaultConf)

		renderConfig := new(defaultFormat, inputFile, outputFile, outputFormat)
		if err := processRender(renderConfig); err != nil {
			er(err)
		} else {
			fmt.Printf("Done rendering secrets to %s\n", outputFile)
		}
	},
}

func createVaultClient(v *VaultConfig) *secret.VaultClient {
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

func processRender(r *RenderConfig) error {
	if r.InputFile == "" {
		return errorInputFileNotFound
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
		return processYamlFile(r)
	case "json":
		return processJSONFile(r)
	}
	return nil
}

func processYamlFile(r *RenderConfig) error {
	d, err := ioutil.ReadFile(r.InputFile)
	if err != nil {
		return err
	}
	flattened := helpers.FlattenYaml(d)
	if flattened == nil {
		return errors.New("error flattening yaml, look into fixing it")
	}

	values, errs := vaultClient.ResolveValues(flattened)
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
	return dirs.WriteContentToFile(r.OutputFile, res, true)
}

func processJSONFile(r *RenderConfig) error {
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

package dirs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rnkoaa/vault-env/helpers"
	"gopkg.in/yaml.v2"
)

var (
	errorFileExists = errors.New("error file exists, not overwriting")
)

func readJSONFile(name string) (map[string]interface{}, error) {
	if !helpers.Exists(name) {
		return nil, fmt.Errorf("file %s does not exist", name)
	}
	d, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(d, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// WriteContentToFile -
func WriteContentToFile(path string, content []byte, overwriteExisting bool) error {
	if helpers.Exists(path) && !overwriteExisting {
		return errorFileExists
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}
	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}

// ReadYamlMapFile -
func ReadYamlMapFile(path string) (map[interface{}]interface{}, error) {
	configData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(configData), &m)
	return m, err
}

// ReadYamlSeqFile -
// Reads a yaml file that is configured as a sequence
func ReadYamlSeqFile(path string) ([]interface{}, error) {
	configData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	m := make([]interface{}, 0)
	err = yaml.Unmarshal([]byte(configData), &m)
	return m, err
}

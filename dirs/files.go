package dirs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/rnkoaa/vault-env/helpers"
	"gopkg.in/yaml.v2"
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

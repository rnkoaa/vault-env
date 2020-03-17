package helpers

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/doublerebel/bellows"
	"github.com/ghodss/yaml"
)

var (
	// ErrorFileNotFound - file not found error
	ErrorFileNotFound = fmt.Errorf("file not found")

	// ErrorFileTypeUndetermined - unable to determine file type
	ErrorFileTypeUndetermined = fmt.Errorf("file type not determined")
)

// func resolveSecretConfigFile(vaultTeam, vaultEnv string, secretConfigFile *domain.ClusterConfig) error {
// 	if secretConfigFile.Content != nil {
// 		flattenedContent := helpers.FlattenYaml(secretConfigFile.Content)
// 		resolvedSecrets, errs := secrets.ResolveSecrets(vaultTeam, vaultEnv, flattenedContent)
// 		if len(errs) > 0 {
// 			var s = fmt.Sprintf("There are %d errors when resolving secrets.\n", len(errs))
// 			for _, err := range errs {
// 				s += fmt.Sprintf("%s\n", err)
// 			}
// 			return fmt.Errorf("%s", s)
// 		}

// 		ymlData, err := helpers.ExpandYaml(resolvedSecrets)
// 		if err != nil {
// 			return err
// 		}
// 		secretConfigFile.Content = ymlData
// 	}
// 	return nil
// }

// DeepCopy -
func DeepCopy(dst, src interface{}) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(w)
	err = enc.Encode(src)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(r)
	return dec.Decode(dst)
}

// Flatten - flattens a json or yaml content into simple kv map
func Flatten(content []byte) map[string]interface{} {
	var b map[string]interface{}
	var res map[string]interface{}
	e := json.Unmarshal(content, &b)
	if e != nil {
		log.Printf("error unmarshalling json to bytes, %#v", e.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			res = make(map[string]interface{})
			//var ok bool
			//err, ok := r.(error)
			//if !ok {
			//	err = fmt.Errorf("pkg: %v", r)
			//}
		}
	}()

	res = bellows.Flatten(b)
	return res
}

// FlattenYaml - flattens a yaml content into simple kv maps
func FlattenYaml(yamlContent []byte) map[string]interface{} {

	yamlBytes, err := yaml.YAMLToJSON(yamlContent)
	if err != nil {
		log.Printf("error unmarshalling content %v", err)
		return nil
	}

	return Flatten(yamlBytes)
}

// ExpandYaml -
// Expands a yaml content from flattened content into a well formed yaml structure
func ExpandYaml(yamlMap map[string]interface{}) ([]byte, error) {
	expandedYamlMap := bellows.Expand(yamlMap)
	b, err := json.Marshal(expandedYamlMap)
	if err != nil {
		log.Printf("error converting map to json")
		return nil, err
	}
	ymlData, err := yaml.JSONToYAML(b)
	if err != nil {
		log.Printf("error converting data into yml")
		return nil, err
	}

	return ymlData, nil
}

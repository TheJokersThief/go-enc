package enc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ENCs        map[string]*ENC
	GlobPattern string
}

func NewConfig(globPatttern string) *Config {
	matchingFiles, err := filepath.Glob(globPatttern)
	err_check(err)

	if matchingFiles == nil {
		panic(errors.New("No files matched that glob pattern."))
	}

	c := &Config{
		ENCs:        make(map[string]*ENC, len(matchingFiles)),
		GlobPattern: globPatttern,
	}

	for _, file := range matchingFiles {
		var enc *ENC
		switch extension := strings.ToLower(filepath.Ext(file)); extension {
		case ".json":
			enc = NewENC("json")
			c.processJSONFile(file, enc)
		case ".yaml", ".yml":
			enc = NewENC("yaml")
			c.processYAMLFile(file, enc)
		default:
			panic(errors.New("Unrecognised file extension, expecting: json|yaml"))
		}

		filename := filepath.Base(file)
		c.ENCs[filename] = enc
	}

	return c
}

func (c *Config) processJSONFile(filepath string, enc *ENC) {
	data, file_err := ioutil.ReadFile(filepath)
	err_check(file_err)

	var raw_enc map[string]interface{}
	json_parse_err := json.Unmarshal(data, &raw_enc)
	err_check(json_parse_err)

	c.processRawENC(raw_enc, enc)
}

func (c *Config) processYAMLFile(filepath string, enc *ENC) {
	data, file_err := ioutil.ReadFile(filepath)
	err_check(file_err)

	var raw_enc map[string]interface{}
	yaml_parse_err := yaml.Unmarshal(data, &raw_enc)
	err_check(yaml_parse_err)

	// YAML unmarshalling returns type map[interface{}]interface{} regardless of provided type
	// so until that's fixed, some conversion has to take place
	for k, v := range raw_enc {
		raw_enc[k] = stringifyYAMLMapKeys(v)
	}

	c.processRawENC(raw_enc, enc)
}

func (c *Config) processRawENC(raw_enc map[string]interface{}, enc *ENC) {
	for nodegroup, attributes := range raw_enc {
		attrs := attributes.(map[string]interface{})

		var (
			parent     string
			classes    map[string]interface{}
			nodes      []string
			parameters map[string]interface{}
			ok         bool
		)

		if parent, ok = attrs["parent"].(string); !ok {
			parent = ""
		}

		if classes, ok = attrs["classes"].(map[string]interface{}); !ok {
			classes = make(map[string]interface{}, 0)
		}

		if nodes, ok = attrs["nodes"].([]string); !ok {
			nodes = make([]string, 0)
		}

		if parameters, ok = attrs["parameters"].(map[string]interface{}); !ok {
			parameters = make(map[string]interface{}, 0)
		}

		enc.AddNodegroup(
			nodegroup,
			parent,
			classes,
			nodes,
			parameters)
	}
}

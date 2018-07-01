package enc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	globalConfig *Config
)

// Config stores the configuration for our ENC
type Config struct {
	ENCs        map[string]*ENC
	GlobPattern string
}

func GetGlobalConfig() *Config {
	return globalConfig
}

// NewConfig generates a new ENC from the config. One ENC for each file matched by the glob pattern
func NewConfig(globPatttern string) *Config {
	matchingFiles, err := filepath.Glob(globPatttern)
	errCheck(err)

	if matchingFiles == nil {
		panic(errors.New("No files matched that glob pattern"))
	}

	c := &Config{
		ENCs:        make(map[string]*ENC, len(matchingFiles)),
		GlobPattern: globPatttern,
	}

	for _, file := range matchingFiles {
		var enc *ENC
		extension := strings.ToLower(filepath.Ext(file))
		switch extension {
		case ".json":
			enc = NewENC("json", file)
			c.processJSONFile(file, enc)
		case ".yaml", ".yml":
			enc = NewENC("yaml", file)
			c.processYAMLFile(file, enc)
		default:
			panic(errors.New("Unrecognised file extension, expecting: json|yaml"))
		}

		filename := filepath.Base(file)
		// Strip extension from filename
		filename = filename[0 : len(filename)-len(extension)]
		c.ENCs[filename] = enc
	}

	globalConfig = c

	return c
}

func (c *Config) WriteOutENC() {
	for _, current_enc := range c.ENCs {
		file, fileErr := os.Create(current_enc.FileName)
		defer file.Close()
		errCheck(fileErr)

		var encContents []byte
		var marshalErr error

		switch extension := strings.ToLower(filepath.Ext(current_enc.FileName)); extension {
		case ".json":
			encContents, marshalErr = json.Marshal(current_enc.Nodegroups)
		case ".yaml", ".yml":
			encContents, marshalErr = yaml.Marshal(current_enc.Nodegroups)
		}
		errCheck(marshalErr)

		file.Write(encContents)
		file.Sync()
	}
}

func (c *Config) processJSONFile(filepath string, enc *ENC) {
	data, fileErr := ioutil.ReadFile(filepath)
	errCheck(fileErr)

	var rawEnc map[string]interface{}
	jsonParseErr := json.Unmarshal(data, &rawEnc)
	errCheck(jsonParseErr)

	c.processRawENC(rawEnc, enc)
}

func (c *Config) processYAMLFile(filepath string, enc *ENC) {
	data, fileErr := ioutil.ReadFile(filepath)
	errCheck(fileErr)

	var rawEnc map[string]interface{}
	yamlParseErr := yaml.Unmarshal(data, &rawEnc)
	errCheck(yamlParseErr)

	// YAML unmarshalling returns type map[interface{}]interface{} regardless of provided type
	// so until that's fixed, some conversion has to take place
	for k, v := range rawEnc {
		rawEnc[k] = stringifyYAMLMapKeys(v)
	}

	c.processRawENC(rawEnc, enc)
}

func (c *Config) processRawENC(rawEnc map[string]interface{}, enc *ENC) {
	nodegroupNodes := make(map[string][]string, 0)

	for nodegroup, attributes := range rawEnc {
		attrs := attributes.(map[string]interface{})

		var (
			parent     string
			classes    map[string]interface{}
			parameters map[string]interface{}
			ok         bool
		)

		if parent, ok = attrs["parent"].(string); !ok {
			parent = ""
		}

		if classes, ok = attrs["classes"].(map[string]interface{}); !ok {
			classes = make(map[string]interface{}, 0)
		}

		if parameters, ok = attrs["parameters"].(map[string]interface{}); !ok {
			parameters = make(map[string]interface{}, 0)
		}

		enc.AddNodegroup(
			nodegroup,
			parent,
			classes,
			make([]string, 0),
			parameters)

		nodegroupNodes[nodegroup] = make([]string, 0)
		if attrs["nodes"] != nil {
			for _, node := range attrs["nodes"].([]interface{}) {
				nodegroupNodes[nodegroup] = append(nodegroupNodes[nodegroup], node.(string))
			}
		}
	}

	for nodegroup, nodes := range nodegroupNodes {
		enc.AddNodes(nodegroup, nodes)
	}
}

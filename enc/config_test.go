package enc

import (
  "os"
  "testing"

  "github.com/derekparker/trie"
  "github.com/stretchr/testify/assert"
)

var yamlData = `
globals:
  classes:
    dns_caching:
      cache_timeout: 50
      size: 200
    extra_archives:
      extra_archives:
        nested_param:
          nested_val1: 1
          nested_val2: 2
        next_nested_param:
          nested_val1: 1
          nested_val2: 2
  parameters:
    test_value: "I'm a global value"
    admin_uid: "1234567"
  parent:
website:
  classes:
    dns_caching:
      cache_timeout: 10
    extra_archives:
      extra_archives:
        nested_param:
          nested_val1: 1
          nested_val2: 2
        last_nested_param:
          nested_val1: 1
          nested_val2: 2
  parameters:
    test_value: "I'm a single value"
  parent: globals
`

var jsonData = `
{
  "globals": {
    "classes": {
      "dns_caching": {
        "cache_timeout": 50,
        "size": 200
      },
      "extra_archives": {
        "extra_archives": {
          "nested_param": {
            "nested_val1": 1,
            "nested_val2": 2
          },
          "next_nested_param": {
            "nested_val1": 1,
            "nested_val2": 2
          }
        }
      }
    },
    "parameters": {
      "test_value": "I'm a global value",
      "admin_uid": "1234567"
    },
    "parent": null
  },
  "website": {
    "classes": {
      "dns_caching": {
        "cache_timeout": 10
      },
      "extra_archives": {
        "extra_archives": {
          "nested_param": {
            "nested_val1": 1,
            "nested_val2": 2
          },
          "last_nested_param": {
            "nested_val1": 1,
            "nested_val2": 2
          }
        }
      }
    },
    "parameters": {
      "test_value": "I'm a single value"
    },
    "parent": "globals"
  }
}
`

var (
  jsonFile = "/tmp/enc_test-json_data.json"
  yamlFile = "/tmp/enc_test-yaml_data.yaml"
)

func init() {
  // JSON setup
  jsonFile, jsonOpenErr := os.Create(jsonFile)
  defer jsonFile.Close()
  if jsonOpenErr != nil {
    panic(jsonOpenErr)
  }

  _, jsonWriteErr := jsonFile.WriteString(jsonData)
  if jsonWriteErr != nil {
    panic(jsonWriteErr)
  }

  jsonFile.Sync()

  // YAML setup
  yamlFile, yamlOpenErr := os.Create(yamlFile)
  defer yamlFile.Close()
  if yamlOpenErr != nil {
    panic(yamlOpenErr)
  }

  _, yamlWriteErr := yamlFile.WriteString(yamlData)
  if yamlWriteErr != nil {
    panic(yamlWriteErr)
  }

  yamlFile.Sync()
}

func TestNewConfig(t *testing.T) {
  testNewJSONConfig(t)
  testNewYAMLConfig(t)
}

func testNewJSONConfig(t *testing.T) {
  assert := assert.New(t)

  wantEnc := ENC{
    Nodegroups: map[string]Nodegroup{
      "globals": Nodegroup{
        Classes: map[string]interface{}{
          "dns_caching": map[string]interface{}{
            "cache_timeout": float64(50),
            "size":          float64(200),
          },
          "extra_archives": map[string]interface{}{
            "extra_archives": map[string]interface{}{
              "nested_param": map[string]interface{}{
                "nested_val1": float64(1),
                "nested_val2": float64(2),
              },
              "next_nested_param": map[string]interface{}{
                "nested_val1": float64(1),
                "nested_val2": float64(2),
              },
            },
          },
        },
        Nodes: []string{},
        Parameters: map[string]interface{}{
          "test_value": "I'm a global value",
          "admin_uid":  "1234567",
        },
        Parent: "",
      },
      "website": Nodegroup{
        Classes: map[string]interface{}{
          "dns_caching": map[string]interface{}{
            "cache_timeout": float64(10),
          },
          "extra_archives": map[string]interface{}{
            "extra_archives": map[string]interface{}{
              "nested_param": map[string]interface{}{
                "nested_val1": float64(1),
                "nested_val2": float64(2),
              },
              "last_nested_param": map[string]interface{}{
                "nested_val1": float64(1),
                "nested_val2": float64(2),
              },
            },
          },
        },
        Nodes: []string{},
        Parameters: map[string]interface{}{
          "test_value": "I'm a single value",
        },
        Parent: "globals",
      },
    },
    Nodes:      trie.New(),
    ConfigType: "json",
  }

  gotJSONConfig := NewConfig(jsonFile)
  assert.Equal(wantEnc, *gotJSONConfig.ENCs["enc_test-json_data.json"])
}

func testNewYAMLConfig(t *testing.T) {
  assert := assert.New(t)

  wantEnc := ENC{
    Nodegroups: map[string]Nodegroup{
      "globals": Nodegroup{
        Classes: map[string]interface{}{
          "dns_caching": map[string]interface{}{
            "cache_timeout": 50,
            "size":          200,
          },
          "extra_archives": map[string]interface{}{
            "extra_archives": map[string]interface{}{
              "nested_param": map[string]interface{}{
                "nested_val1": 1,
                "nested_val2": 2,
              },
              "next_nested_param": map[string]interface{}{
                "nested_val1": 1,
                "nested_val2": 2,
              },
            },
          },
        },
        Nodes: []string{},
        Parameters: map[string]interface{}{
          "test_value": "I'm a global value",
          "admin_uid":  "1234567",
        },
        Parent: "",
      },
      "website": Nodegroup{
        Classes: map[string]interface{}{
          "dns_caching": map[string]interface{}{
            "cache_timeout": 10,
          },
          "extra_archives": map[string]interface{}{
            "extra_archives": map[string]interface{}{
              "nested_param": map[string]interface{}{
                "nested_val1": 1,
                "nested_val2": 2,
              },
              "last_nested_param": map[string]interface{}{
                "nested_val1": 1,
                "nested_val2": 2,
              },
            },
          },
        },
        Nodes: []string{},
        Parameters: map[string]interface{}{
          "test_value": "I'm a single value",
        },
        Parent: "globals",
      },
    },
    Nodes:      trie.New(),
    ConfigType: "yaml",
  }

  gotYAMLConfig := NewConfig(yamlFile)
  assert.Equal(wantEnc, *gotYAMLConfig.ENCs["enc_test-yaml_data.yaml"])
}

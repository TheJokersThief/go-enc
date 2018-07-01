package enc

import (
	"testing"

	"github.com/derekparker/trie"
	"github.com/stretchr/testify/assert"
)

var (
	conf *Config
)

func init() {
	conf = NewConfig("/tmp/enc_test-json_data.json")
}

func TestNewENC(t *testing.T) {
	assert := assert.New(t)
	nodesTrie := trie.New()
	want := ENC{
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{},
		Nodes:      nodesTrie,
		ConfigType: "json",
	}

	got := NewENC("json", "/tmp/enc_test-json_data.json")
	assert.Equal(want, *got)
}

func TestAddNodegroup(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent:      "",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup": wantNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotNodegroup, nodegroupErr := gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))

	assert.Equal(wantEnc, *gotEnc)
	assert.Equal(wantNodegroup, *gotNodegroup)
	assert.Nil(nodegroupErr)
}

func TestRemoveNodegroup(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent:      "",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotNodegroup, nodegroupErr := gotEnc.RemoveNodegroup("wantNodegroup")

	assert.Equal(wantEnc, *gotEnc)
	assert.Equal(wantNodegroup, *gotNodegroup)
	assert.Nil(nodegroupErr)
}

func TestGetNodegroup(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent:      "",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup": wantNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))

	gotNodegroup, nodegroupErr := gotEnc.GetNodegroup("wantNodegroup")

	assert.Equal(wantEnc, *gotEnc)
	assert.Equal(wantNodegroup, *gotNodegroup)
	assert.Nil(nodegroupErr)
}

func TestAddNode(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent:  "",
		Classes: make(map[string]interface{}, 0),
		Nodes: []string{
			"node-0001",
		},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup": wantNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))

	gotNodegroup, nodegroupErr := gotEnc.AddNode("wantNodegroup", "node-0001")

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes

	assert.Equal(wantEnc, *gotEnc)
	assert.Equal(wantNodegroup, *gotNodegroup)
	assert.Nil(nodegroupErr)
}

func TestAddNodes(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent:  "",
		Classes: make(map[string]interface{}, 0),
		Nodes: []string{
			"node-0001",
			"node-0002",
		},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup": wantNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))

	gotNodegroup, nodegroupErr := gotEnc.AddNodes("wantNodegroup", []string{"node-0001", "node-0002"})

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes

	assert.Equal(wantEnc, *gotEnc)
	assert.Equal(wantNodegroup, *gotNodegroup)
	assert.Nil(nodegroupErr)
}

func TestGetParentChain(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent:      "",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	subNodegroup := Nodegroup{
		Parent:      "wantNodegroup@enc_test-json_data",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	subSubNodegroup := Nodegroup{
		Parent:      "subNodegroup@enc_test-json_data",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup":   wantNodegroup,
			"subNodegroup":    subNodegroup,
			"subSubNodegroup": subSubNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	wantChain := []string{"subSubNodegroup@enc_test-json_data", "subNodegroup@enc_test-json_data", "wantNodegroup@enc_test-json_data"}

	conf.ENCs["enc_test-json_data"].Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddNodegroup("subNodegroup", "wantNodegroup@enc_test-json_data", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddNodegroup("subSubNodegroup", "subNodegroup@enc_test-json_data", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))

	gotChain := gotEnc.getParentChain("subSubNodegroup")

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes

	assert.Equal(wantEnc, *gotEnc)
	assert.Equal(wantChain, gotChain)

}

func TestGetLongestChain(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent:  "",
		Classes: make(map[string]interface{}, 0),
		Nodes: []string{
			"node-0001",
		},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	subNodegroup := Nodegroup{
		Parent:  "wantNodegroup@enc_test-json_data",
		Classes: make(map[string]interface{}, 0),
		Nodes: []string{
			"node-0001",
		},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	subSubNodegroup := Nodegroup{
		Parent:      "subNodegroup@enc_test-json_data",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup":   wantNodegroup,
			"subNodegroup":    subNodegroup,
			"subSubNodegroup": subSubNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.Name = "enc_test-json_data"
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddNodegroup("subNodegroup", "wantNodegroup@enc_test-json_data", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddNodegroup("subSubNodegroup", "subNodegroup@enc_test-json_data", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddNode("wantNodegroup", "node-0001")
	gotEnc.AddNode("subNodegroup", "node-0001")

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes

	assert.Equal("node-0001-wantNodegroup@enc_test-json_data-subNodegroup@enc_test-json_data", gotEnc.getLongestChain("node-0001"))
	assert.Equal(wantEnc, *gotEnc)
}

func TestRemoveNode(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent:  "",
		Classes: make(map[string]interface{}, 0),
		Nodes: []string{
			"node-0001",
		},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	subNodegroup := Nodegroup{
		Parent:  "wantNodegroup@enc_test-json_data",
		Classes: make(map[string]interface{}, 0),
		Nodes: []string{
			"node-0001",
		},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	subSubNodegroup := Nodegroup{
		Parent:      "subNodegroup@enc_test-json_data",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup":   wantNodegroup,
			"subNodegroup":    subNodegroup,
			"subSubNodegroup": subSubNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.Name = "enc_test-json_data"
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddNodegroup("subNodegroup", "wantNodegroup@enc_test-json_data", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddNodegroup("subSubNodegroup", "subNodegroup@enc_test-json_data", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddNode("wantNodegroup", "node-0001")
	gotEnc.AddNode("subNodegroup", "node-0001")

	// Because of the nature of a trie, we only want the key to disappear if
	// it has no children. wantNodegroup has one child so should remain
	gotEnc.RemoveNode("wantNodegroup", "node-0001")
	assert.NotNil(gotEnc.Nodes.Find("node-0001-wantNodegroup@enc_test-json_data"))

	// subNodegroup on the other hand has no children so will be removed
	gotEnc.RemoveNode("subNodegroup", "node-0001")
	assert.Nil(gotEnc.Nodes.Find("node-0001-wantNodegroup@enc_test-json_data-subNodegroup@enc_test-json_data"))

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes
	assert.Equal(wantEnc, *gotEnc)
}

func TestAddParameter(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent:  "",
		Classes: make(map[string]interface{}, 0),
		Nodes:   []string{},
		Parameters: map[string]interface{}{
			"test_param": "test_value",
		},
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup": wantNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddParameter("wantNodegroup", "test_param", "test_value")

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes

	assert.Equal(wantEnc, *gotEnc)
}

func TestRemoveParameter(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent:      "",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup": wantNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddParameter("wantNodegroup", "test_param", "test_value")
	gotEnc.RemoveParameter("wantNodegroup", "test_param")

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes

	assert.Equal(wantEnc, *gotEnc)
}

func TestAddClass(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent: "",
		Classes: map[string]interface{}{
			"test_class": map[string]interface{}{},
		},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup": wantNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddClass("wantNodegroup", "test_class")

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes

	assert.Equal(wantEnc, *gotEnc)
}

func TestRemoveClass(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent:      "",
		Classes:     map[string]interface{}{},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup": wantNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddClass("wantNodegroup", "test_class")
	gotEnc.RemoveClass("wantNodegroup", "test_class")

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes

	assert.Equal(wantEnc, *gotEnc)
}

func TestAddClassParameter(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent: "",
		Classes: map[string]interface{}{
			"test_class": map[string]interface{}{
				"test_class_param": "test_value",
			},
		},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup": wantNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddClass("wantNodegroup", "test_class")
	gotEnc.AddClassParameter("wantNodegroup", "test_class", "test_class_param", "test_value")

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes

	assert.Equal(wantEnc, *gotEnc)
}

func TestRemoveClassParameter(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent: "",
		Classes: map[string]interface{}{
			"test_class": map[string]interface{}{},
		},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup": wantNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddClass("wantNodegroup", "test_class")
	gotEnc.AddClassParameter("wantNodegroup", "test_class", "test_class_param", "test_value")
	gotEnc.RemoveClassParameter("wantNodegroup", "test_class", "test_class_param")

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes

	assert.Equal(wantEnc, *gotEnc)
}

func TestSetParent(t *testing.T) {
	assert := assert.New(t)

	parentNodegroup := Nodegroup{
		Parent:      "",
		Classes:     map[string]interface{}{},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	wantNodegroup := Nodegroup{
		Parent:      "test_parent",
		Classes:     map[string]interface{}{},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"test_parent":   parentNodegroup,
			"wantNodegroup": wantNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotEnc.AddNodegroup("test_parent", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotNodegroup, gotErr := gotEnc.SetParent("wantNodegroup", "test_parent")

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes

	assert.Nil(gotErr)
	assert.Equal(wantNodegroup, *gotNodegroup)
	assert.Equal(wantEnc, *gotEnc)
}

func TestSetEnvironment(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent:      "",
		Classes:     map[string]interface{}{},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "test_env",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup": wantNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	gotNodegroup, gotErr := gotEnc.SetEnvironment("wantNodegroup", "test_env")

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes

	assert.Nil(gotErr)
	assert.Equal(wantNodegroup, *gotNodegroup)
	assert.Equal(wantEnc, *gotEnc)
}

func TestGetNode(t *testing.T) {
	assert := assert.New(t)

	wantNodegroup := Nodegroup{
		Parent: "",
		Classes: map[string]interface{}{
			"test_class": map[string]interface{}{
				"override_me": "I should not be here",
			},
			"third_class": map[string]interface{}{
				"unique_test": "I've never been overriden",
			},
		},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "env_one",
	}

	subNodegroup := Nodegroup{
		Parent: "wantNodegroup",
		Classes: map[string]interface{}{
			"test_class": map[string]interface{}{
				"unique_test": "I've never been overriden",
			},
			"second_class": map[string]interface{}{
				"unique_test": "I've never been overriden",
			},
		},
		Nodes: []string{},
		Parameters: map[string]interface{}{
			"test_param": "test_value",
		},
		Environment: "env_two",
	}

	subSubNodegroup := Nodegroup{
		Parent: "subNodegroup",
		Classes: map[string]interface{}{
			"test_class": map[string]interface{}{
				"override_me": "I'm legit",
			},
			"fourth_class": map[string]interface{}{
				"unique_test": "I've never been overriden",
			},
		},
		Nodes: []string{
			"node-0001",
		},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	wantEnc := ENC{
		ConfigLink: conf,
		Name:       "enc_test-json_data",
		FileName:   "/tmp/enc_test-json_data.json",
		Nodegroups: map[string]Nodegroup{
			"wantNodegroup":   wantNodegroup,
			"subNodegroup":    subNodegroup,
			"subSubNodegroup": subSubNodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	wantNode := Nodegroup{
		Parent: "subNodegroup",
		Classes: map[string]interface{}{
			"test_class": map[string]interface{}{
				"override_me": "I'm legit",
				"unique_test": "I've never been overriden",
			},
			"second_class": map[string]interface{}{
				"unique_test": "I've never been overriden",
			},
			"third_class": map[string]interface{}{
				"unique_test": "I've never been overriden",
			},
			"fourth_class": map[string]interface{}{
				"unique_test": "I've never been overriden",
			},
		},
		Nodes: []string{
			"node-0001",
		},
		Parameters: map[string]interface{}{
			"test_param": "test_value",
		},
		Environment: "env_two",
	}

	gotEnc := conf.ENCs["enc_test-json_data"]
	gotEnc.Nodegroups = make(map[string]Nodegroup, 0)
	gotEnc.AddNodegroup("wantNodegroup", "", make(map[string]interface{}), []string{}, make(map[string]interface{}))
	gotEnc.AddNodegroup("subNodegroup", "wantNodegroup", make(map[string]interface{}), []string{}, make(map[string]interface{}))
	gotEnc.AddNodegroup("subSubNodegroup", "subNodegroup", make(map[string]interface{}), []string{}, make(map[string]interface{}))

	gotEnc.AddNode("subSubNodegroup", "node-0001")

	gotEnc.SetEnvironment("wantNodegroup", "env_one")
	gotEnc.SetEnvironment("subNodegroup", "env_two")

	gotEnc.AddClass("wantNodegroup", "test_class")
	gotEnc.AddClass("wantNodegroup", "third_class")
	gotEnc.AddClass("subNodegroup", "test_class")
	gotEnc.AddClass("subNodegroup", "second_class")
	gotEnc.AddClass("subSubNodegroup", "test_class")
	gotEnc.AddClass("subSubNodegroup", "fourth_class")

	gotEnc.AddParameter("subNodegroup", "test_param", "test_value")

	gotEnc.AddClassParameter("wantNodegroup", "test_class", "override_me", "I should not be here")
	gotEnc.AddClassParameter("wantNodegroup", "third_class", "unique_test", "I've never been overriden")

	gotEnc.AddClassParameter("subNodegroup", "test_class", "unique_test", "I've never been overriden")
	gotEnc.AddClassParameter("subNodegroup", "second_class", "unique_test", "I've never been overriden")

	gotEnc.AddClassParameter("subSubNodegroup", "test_class", "override_me", "I'm legit")
	gotEnc.AddClassParameter("subSubNodegroup", "fourth_class", "unique_test", "I've never been overriden")

	gotNode, gotErr := gotEnc.GetNode("node-0001")

	// We aren't testing the trie package
	wantEnc.Nodes = gotEnc.Nodes

	assert.Nil(gotErr)
	assert.Equal(wantNode, *gotNode)
}

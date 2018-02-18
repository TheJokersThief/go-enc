package enc

import (
	"testing"

	"github.com/derekparker/trie"
	"github.com/stretchr/testify/assert"
)

func TestNewENC(t *testing.T) {
	assert := assert.New(t)
	nodes_trie := trie.New()
	want := ENC{
		Nodegroups: map[string]ENCNodegroup{},
		Nodes:      nodes_trie,
		ConfigType: "json",
	}

	got := NewENC("json")
	assert.Equal(want, *got)
}

func TestAddNodegroup(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent:      "",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup": want_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	empty_generic_map := make(map[string]interface{}, 0)
	got_nodegroup, nodegroup_err := got_enc.AddNodegroup("want_nodegroup", "", empty_generic_map, []string{}, empty_generic_map)

	assert.Equal(want_enc, *got_enc)
	assert.Equal(want_nodegroup, *got_nodegroup)
	assert.Nil(nodegroup_err)
}

func TestRemoveNodegroup(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent:      "",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	empty_generic_map := make(map[string]interface{}, 0)
	got_enc.AddNodegroup("want_nodegroup", "", empty_generic_map, []string{}, empty_generic_map)
	got_nodegroup, nodegroup_err := got_enc.RemoveNodegroup("want_nodegroup")

	assert.Equal(want_enc, *got_enc)
	assert.Equal(want_nodegroup, *got_nodegroup)
	assert.Nil(nodegroup_err)
}

func TestGetNodegroup(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent:      "",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup": want_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	empty_generic_map := make(map[string]interface{}, 0)
	got_enc.AddNodegroup("want_nodegroup", "", empty_generic_map, []string{}, empty_generic_map)

	got_nodegroup, nodegroup_err := got_enc.GetNodegroup("want_nodegroup")

	assert.Equal(want_enc, *got_enc)
	assert.Equal(want_nodegroup, *got_nodegroup)
	assert.Nil(nodegroup_err)
}

func TestAddNode(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent:  "",
		Classes: make(map[string]interface{}, 0),
		Nodes: []string{
			"node-0001",
		},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup": want_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	empty_generic_map := make(map[string]interface{}, 0)
	got_enc.AddNodegroup("want_nodegroup", "", empty_generic_map, []string{}, empty_generic_map)

	got_nodegroup, nodegroup_err := got_enc.AddNode("want_nodegroup", "node-0001")

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Equal(want_enc, *got_enc)
	assert.Equal(want_nodegroup, *got_nodegroup)
	assert.Nil(nodegroup_err)
}

func TestAddNodes(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent:  "",
		Classes: make(map[string]interface{}, 0),
		Nodes: []string{
			"node-0001",
			"node-0002",
		},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup": want_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	empty_generic_map := make(map[string]interface{}, 0)
	got_enc.AddNodegroup("want_nodegroup", "", empty_generic_map, []string{}, empty_generic_map)

	got_nodegroup, nodegroup_err := got_enc.AddNodes("want_nodegroup", []string{"node-0001", "node-0002"})

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Equal(want_enc, *got_enc)
	assert.Equal(want_nodegroup, *got_nodegroup)
	assert.Nil(nodegroup_err)
}

func TestGetParentChain(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent:      "",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	sub_nodegroup := ENCNodegroup{
		Parent:      "want_nodegroup",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	sub_sub_nodegroup := ENCNodegroup{
		Parent:      "sub_nodegroup",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup":    want_nodegroup,
			"sub_nodegroup":     sub_nodegroup,
			"sub_sub_nodegroup": sub_sub_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	want_chain := []string{"sub_sub_nodegroup", "sub_nodegroup", "want_nodegroup"}

	got_enc := NewENC("json")
	empty_generic_map := make(map[string]interface{}, 0)
	got_enc.AddNodegroup("want_nodegroup", "", empty_generic_map, []string{}, empty_generic_map)
	got_enc.AddNodegroup("sub_nodegroup", "want_nodegroup", empty_generic_map, []string{}, empty_generic_map)
	got_enc.AddNodegroup("sub_sub_nodegroup", "sub_nodegroup", empty_generic_map, []string{}, empty_generic_map)

	got_chain := got_enc.getParentChain("sub_sub_nodegroup")

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Equal(want_enc, *got_enc)
	assert.Equal(want_chain, got_chain)

}

func TestGetLongestChain(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent:  "",
		Classes: make(map[string]interface{}, 0),
		Nodes: []string{
			"node-0001",
		},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	sub_nodegroup := ENCNodegroup{
		Parent:  "want_nodegroup",
		Classes: make(map[string]interface{}, 0),
		Nodes: []string{
			"node-0001",
		},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	sub_sub_nodegroup := ENCNodegroup{
		Parent:      "sub_nodegroup",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup":    want_nodegroup,
			"sub_nodegroup":     sub_nodegroup,
			"sub_sub_nodegroup": sub_sub_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	empty_generic_map := make(map[string]interface{}, 0)
	got_enc.AddNodegroup("want_nodegroup", "", empty_generic_map, []string{}, empty_generic_map)
	got_enc.AddNodegroup("sub_nodegroup", "want_nodegroup", empty_generic_map, []string{}, empty_generic_map)
	got_enc.AddNodegroup("sub_sub_nodegroup", "sub_nodegroup", empty_generic_map, []string{}, empty_generic_map)
	got_enc.AddNode("want_nodegroup", "node-0001")
	got_enc.AddNode("sub_nodegroup", "node-0001")

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Equal("node-0001-sub_nodegroup-want_nodegroup", got_enc.getLongestChain("node-0001"))
	assert.Equal(want_enc, *got_enc)
}

func TestRemoveNode(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent:  "",
		Classes: make(map[string]interface{}, 0),
		Nodes: []string{
			"node-0001",
		},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	sub_nodegroup := ENCNodegroup{
		Parent:  "want_nodegroup",
		Classes: make(map[string]interface{}, 0),
		Nodes: []string{
			"node-0001",
		},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	sub_sub_nodegroup := ENCNodegroup{
		Parent:      "sub_nodegroup",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  make(map[string]interface{}, 0),
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup":    want_nodegroup,
			"sub_nodegroup":     sub_nodegroup,
			"sub_sub_nodegroup": sub_sub_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	empty_generic_map := make(map[string]interface{}, 0)
	got_enc.AddNodegroup("want_nodegroup", "", empty_generic_map, []string{}, empty_generic_map)
	got_enc.AddNodegroup("sub_nodegroup", "want_nodegroup", empty_generic_map, []string{}, empty_generic_map)
	got_enc.AddNodegroup("sub_sub_nodegroup", "sub_nodegroup", empty_generic_map, []string{}, empty_generic_map)
	got_enc.AddNode("want_nodegroup", "node-0001")
	got_enc.AddNode("sub_nodegroup", "node-0001")

	got_enc.RemoveNode("want_nodegroup", "node-0001")

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Equal(want_enc, *got_enc)
	assert.Nil(got_enc.Nodes.Find("node-0001-want_nodegroup"))
	assert.NotNil(got_enc.Nodes.Find("node-0001-sub_nodegroup-want_nodegroup"))
}

func TestAddParameter(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent:  "",
		Classes: make(map[string]interface{}, 0),
		Nodes:   []string{},
		Parameters: map[string]interface{}{
			"test_param": "test_value",
		},
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup": want_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	got_enc.AddNodegroup("want_nodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	got_enc.AddParameter("want_nodegroup", "test_param", "test_value")

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Equal(want_enc, *got_enc)
}

func TestRemoveParameter(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent:      "",
		Classes:     make(map[string]interface{}, 0),
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup": want_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	got_enc.AddNodegroup("want_nodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	got_enc.AddParameter("want_nodegroup", "test_param", "test_value")
	got_enc.RemoveParameter("want_nodegroup", "test_param")

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Equal(want_enc, *got_enc)
}

func TestAddClass(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent: "",
		Classes: map[string]interface{}{
			"test_class": map[string]interface{}{},
		},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup": want_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	got_enc.AddNodegroup("want_nodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	got_enc.AddClass("want_nodegroup", "test_class")

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Equal(want_enc, *got_enc)
}

func TestRemoveClass(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent:      "",
		Classes:     map[string]interface{}{},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup": want_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	got_enc.AddNodegroup("want_nodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	got_enc.AddClass("want_nodegroup", "test_class")
	got_enc.RemoveClass("want_nodegroup", "test_class")

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Equal(want_enc, *got_enc)
}

func TestAddClassParameter(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
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

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup": want_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	got_enc.AddNodegroup("want_nodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	got_enc.AddClass("want_nodegroup", "test_class")
	got_enc.AddClassParameter("want_nodegroup", "test_class", "test_class_param", "test_value")

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Equal(want_enc, *got_enc)
}

func TestRemoveClassParameter(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent: "",
		Classes: map[string]interface{}{
			"test_class": map[string]interface{}{},
		},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup": want_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	got_enc.AddNodegroup("want_nodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	got_enc.AddClass("want_nodegroup", "test_class")
	got_enc.AddClassParameter("want_nodegroup", "test_class", "test_class_param", "test_value")
	got_enc.RemoveClassParameter("want_nodegroup", "test_class", "test_class_param")

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Equal(want_enc, *got_enc)
}

func TestSetParent(t *testing.T) {
	assert := assert.New(t)

	parent_nodegroup := ENCNodegroup{
		Parent:      "",
		Classes:     map[string]interface{}{},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	want_nodegroup := ENCNodegroup{
		Parent:      "test_parent",
		Classes:     map[string]interface{}{},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"test_parent":    parent_nodegroup,
			"want_nodegroup": want_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	got_enc.AddNodegroup("want_nodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	got_enc.AddNodegroup("test_parent", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	got_nodegroup, got_err := got_enc.SetParent("want_nodegroup", "test_parent")

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Nil(got_err)
	assert.Equal(want_nodegroup, *got_nodegroup)
	assert.Equal(want_enc, *got_enc)
}

func TestSetEnvironment(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
		Parent:      "",
		Classes:     map[string]interface{}{},
		Nodes:       []string{},
		Parameters:  map[string]interface{}{},
		Environment: "test_env",
	}

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup": want_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	got_enc := NewENC("json")
	got_enc.AddNodegroup("want_nodegroup", "", make(map[string]interface{}, 0), []string{}, make(map[string]interface{}, 0))
	got_nodegroup, got_err := got_enc.SetEnvironment("want_nodegroup", "test_env")

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Nil(got_err)
	assert.Equal(want_nodegroup, *got_nodegroup)
	assert.Equal(want_enc, *got_enc)
}

func TestGetNode(t *testing.T) {
	assert := assert.New(t)

	want_nodegroup := ENCNodegroup{
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

	sub_nodegroup := ENCNodegroup{
		Parent: "want_nodegroup",
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

	sub_sub_nodegroup := ENCNodegroup{
		Parent: "sub_nodegroup",
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

	want_enc := ENC{
		Nodegroups: map[string]ENCNodegroup{
			"want_nodegroup":    want_nodegroup,
			"sub_nodegroup":     sub_nodegroup,
			"sub_sub_nodegroup": sub_sub_nodegroup,
		},
		Nodes:      trie.New(),
		ConfigType: "json",
	}

	want_node := ENCNodegroup{
		Parent: "sub_nodegroup",
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

	got_enc := NewENC("json")
	got_enc.AddNodegroup("want_nodegroup", "", make(map[string]interface{}), []string{}, make(map[string]interface{}))
	got_enc.AddNodegroup("sub_nodegroup", "want_nodegroup", make(map[string]interface{}), []string{}, make(map[string]interface{}))
	got_enc.AddNodegroup("sub_sub_nodegroup", "sub_nodegroup", make(map[string]interface{}), []string{}, make(map[string]interface{}))

	got_enc.AddNode("sub_sub_nodegroup", "node-0001")

	got_enc.SetEnvironment("want_nodegroup", "env_one")
	got_enc.SetEnvironment("sub_nodegroup", "env_two")

	got_enc.AddClass("want_nodegroup", "test_class")
	got_enc.AddClass("want_nodegroup", "third_class")
	got_enc.AddClass("sub_nodegroup", "test_class")
	got_enc.AddClass("sub_nodegroup", "second_class")
	got_enc.AddClass("sub_sub_nodegroup", "test_class")
	got_enc.AddClass("sub_sub_nodegroup", "fourth_class")

	got_enc.AddParameter("sub_nodegroup", "test_param", "test_value")

	got_enc.AddClassParameter("want_nodegroup", "test_class", "override_me", "I should not be here")
	got_enc.AddClassParameter("want_nodegroup", "third_class", "unique_test", "I've never been overriden")

	got_enc.AddClassParameter("sub_nodegroup", "test_class", "unique_test", "I've never been overriden")
	got_enc.AddClassParameter("sub_nodegroup", "second_class", "unique_test", "I've never been overriden")

	got_enc.AddClassParameter("sub_sub_nodegroup", "test_class", "override_me", "I'm legit")
	got_enc.AddClassParameter("sub_sub_nodegroup", "fourth_class", "unique_test", "I've never been overriden")

	got_node, got_err := got_enc.GetNode("node-0001")

	// We aren't testing the trie package
	want_enc.Nodes = got_enc.Nodes

	assert.Nil(got_err)
	assert.Equal(want_node, *got_node)
}

package cli

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("go-enc", "CLI for interacting with YAML/JSON External Node Classifiers")

	nodegroup       = app.Command("nodegroup", "Actions to do with nodegroups")
	nodegroupAction = nodegroup.Arg("action", "add|remove|get").Required().String()
	nodegroupName   = nodegroup.Arg("nodegroup", "Nodegoup name").Required().String()

	node          = app.Command("node", "Actions to do with single node")
	nodeAction    = node.Arg("action", "add|remove|get").Required().String()
	nodeNodegroup = node.Arg("nodegroup", "Nodegoup name").Required().String()
	nodeNode      = node.Arg("node", "Node").Required().String()
	nodeOutput    = node.Flag("output", "Output format: json|yaml").Default("yaml").Short('o').String()

	nodes          = app.Command("nodes", "Actions to do with single node")
	nodesAdd       = nodes.Arg("add", "Add space-separated list of nodes").Required().String()
	nodesNodegroup = nodes.Arg("nodegroup", "Nodegoup name").Required().String()
	nodesNodes     = StringList(nodes.Arg("nodes", "Nodes").Required())
	nodesOutput    = nodes.Flag("output", "Output format: json|yaml").Default("yaml").Short('o').String()

	param          = app.Command("param", "Actions for parameters")
	paramAction    = param.Arg("action", "add|set|remove").Required().String()
	paramNodegroup = param.Arg("nodegroup", "Nodegoup name").Required().String()
	paramName      = param.Arg("param_name", "Parameter name").Required().String()
	paramValue     = param.Arg("param_value", "Parameter value").Required().String()

	class          = app.Command("class", "Actions for classes")
	classAction    = class.Arg("action", "add|remove").Required().String()
	classNodegroup = class.Arg("nodegroup", "Nodegoup name").Required().String()

	classParam          = app.Command("class_param", "Actions for parameters")
	classParamAction    = classParam.Arg("action", "add|set|remove").Required().String()
	classParamNodegroup = classParam.Arg("nodegroup", "Nodegoup name").Required().String()
	classParamClass     = classParam.Arg("class_name", "Class name").Required().String()
	classParamName      = classParam.Arg("param_name", "Parameter name").Required().String()
	classParamValue     = classParam.Arg("param_value", "Parameter value").Required().String()

	parent    = app.Command("parent", "Set the parent value")
	parentVal = parent.Arg("new_parent", "The new parent value (can be \"\" for none)").Required().String()

	environment    = app.Command("environment", "Set the environment value")
	environmentVal = environment.Arg("new_environment", "The new environment value (can be \"\" for none)").Required().String()
)

func NewCLI() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	// switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// case nodegroupAdd.FullCommand():
	// case nodegroupRemove.FullCommand():
	// }
}

// https://github.com/alecthomas/kingpin

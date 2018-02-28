package cli

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/thejokersthief/go-enc/enc"
)

var (
	app = kingpin.New("go-enc", "CLI for interacting with YAML/JSON External Node Classifiers")

	enc_glob = app.Flag("enc_glob", "Glob pattern for matching ENC files").Default("./*.yaml").Short('g').String()
	enc_name = app.Flag("enc_name", "Name of the ENC you want to perform actions on").Default("production").Short('e').String()

	nodegroup       = app.Command("nodegroup", "Actions to do with nodegroups")
	nodegroupAction = nodegroup.Arg("action", "add|remove|get").Required().String()
	nodegroupName   = nodegroup.Arg("nodegroup", "Nodegoup name").Required().String()
	nodegroupParent = nodegroup.Flag("parent", "Nodegoup parent").Default("").String()

	node          = app.Command("node", "Actions to do with single node")
	nodeAction    = node.Arg("action", "add|remove|get").Required().String()
	nodeNodegroup = node.Arg("nodegroup", "Nodegoup name").Required().String()
	nodeNode      = node.Arg("node", "Node").Required().String()
	nodeOutput    = node.Flag("output", "Output format: json|yaml").Default("yaml").Short('o').String()

	nodes          = app.Command("nodes", "Actions to do with single node")
	nodesAdd       = nodes.Arg("add", "add").Required().String()
	nodesNodegroup = nodes.Arg("nodegroup", "Nodegoup name").Required().String()
	nodesNodes     = StringList(nodes.Arg("nodes", "Space-separated list of nodes").Required())
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
	kingpin.CommandLine.HelpFlag.Short('h')

	config := enc.NewConfig(*enc_glob)
	working_enc, ok := config.ENCs[*enc_name]
	if !ok {
		handleErr(fmt.Errorf("Chosen ENC doesn't exist: %s", enc_name))
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case nodegroup.FullCommand():
		nodegroupCommand(working_enc)
	case node.FullCommand():
		nodeCommand(working_enc)
	case nodes.FullCommand():
		nodesCommand(working_enc)
	case param.FullCommand():
		paramCommand(working_enc)
	case class.FullCommand():
		classCommand(working_enc)
	case classParam.FullCommand():
		classParamCommand(working_enc)
	case parent.FullCommand():
		parentCommand(working_enc)
	case environment.FullCommand():
		environmentCommand(working_enc)
	}
}

func nodegroupCommand(working_enc *enc.ENC) {
	switch *nodegroupAction {
	case "add":
		working_enc.AddNodegroup(*nodegroupName, *nodegroupParent, map[string]interface{}{}, []string{}, map[string]interface{}{})
	case "remove":
		working_enc.RemoveNodegroup(*nodegroupName)
	case "get":
		working_enc.GetNodegroup(*nodegroupName)
	default:
		handleErr(fmt.Errorf("Invalid action for command: [command: %s ; action: %s]", nodegroup.FullCommand(), *nodegroupAction))
	}
}

func nodeCommand(working_enc *enc.ENC) {
	switch *nodeAction {
	case "add":
		working_enc.AddNode(*nodeNodegroup, *nodeNode)
	case "get":
		working_enc.GetNode(*nodeNode)
	case "remove":
		working_enc.RemoveNode(*nodeNodegroup, *nodeNode)
	default:
		handleErr(fmt.Errorf("Invalid action for command: [command: %s ; action: %s]", node.FullCommand(), *nodeAction))
	}
}

func nodesCommand(working_enc *enc.ENC) {
	working_enc.AddNodes(*nodesNodegroup, *nodesNodes)
}

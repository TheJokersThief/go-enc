# Go-ENC
## An External Node Classifier tool
ENCs are used to track configuration across your infrastructure in a single unified location. For more information, [Puppet has a great docs page](https://puppet.com/docs/puppet/5.4/nodes_external.html) on what an ENC is and how to integrate it.

### Features
* Easy interaction/changes to a central config file
* Merged classes and parameters based on parents
* Grouped nodes by nodegroups
* Choice of config format: JSON or YAML
* Command-Line Interface

## Usage
Pick up the latest release binary for your system and try running with the help flag for 
the available commands.

Example output:

```
$ ./go-enc --help
usage: go-enc [<flags>] <command> [<args> ...]

CLI for interacting with YAML/JSON External Node Classifiers

Flags:
      --help                   Show context-sensitive help (also try --help-long and --help-man).
  -g, --enc_glob="./*.yaml"    Glob pattern for matching ENC files
  -e, --enc_name="production"  Name of the ENC you want to perform actions on

Commands:
  help [<command>...]
    Show help.

  nodegroup [<flags>] <action> <nodegroup>
    Actions to do with nodegroups

  node [<flags>] <action> <nodegroup> <node>
    Actions to do with single node

  nodes [<flags>] <add> <nodegroup> <nodes>...
    Actions to do with single node

  param <action> <nodegroup> <param_name> <param_value>
    Actions for parameters

  class <action> <nodegroup> <classname>
    Actions for classes

  class_param <action> <nodegroup> <class_name> <param_name> <param_value>
    Actions for parameters

  parent <nodegroup> <new_parent>
    Set the parent value

  environment <nodegroup> <new_environment>
    Set the environment value
```

### Command Help
You can also pass the help flag after adding your command.

Example output:

```
$ ./go-enc nodegroup --help
usage: go-enc nodegroup [<flags>] <action> <nodegroup>

Actions to do with nodegroups

Flags:
      --help                   Show context-sensitive help (also try --help-long and --help-man).
  -g, --enc_glob="./*.yaml"    Glob pattern for matching ENC files
  -e, --enc_name="production"  Name of the ENC you want to perform actions on
      --parent=""              Nodegoup parent

Args:
  <action>     add|remove|get
  <nodegroup>  Nodegoup name
```


## Development
Go-ENC uses [dep](https://github.com/golang/dep) to manage dependencies.

```
go get github.com/TheJokersThief/go-enc
cd "$GOPATH/src/github.com/TheJokersThief/go-enc"
dep ensure
```

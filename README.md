# Go-ENC
## An External Node Classifier tool
ENCs are used to track configuration across your infrastructure in a single unified location. For more information, [Puppet has a great docs page](https://puppet.com/docs/puppet/5.4/nodes_external.html) on what an ENC is and how to integrate it.

### Features
* Easy interaction/changes to a central config file
* Merged classes and parameters based on parents
* Grouped nodes by nodegroups
* Choice of config format: JSON or YAML
* (WIP) ~~Command-Line Interface~~
* (WIP) ~~REST-API webservice~~

## Development
Go-ENC uses [dep](https://github.com/golang/dep) to manage dependencies.

```
go get github.com/TheJokersThief/go-enc
cd "$GOPATH/src/github.com/TheJokersThief/go-enc"
dep ensure
```

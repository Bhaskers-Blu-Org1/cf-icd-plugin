# cf_icd_plugin
Sample CF plugin for POST of cloud foundry commands to registered webhook broker

Install golang https://golang.org

export GOPATH=/this/project/path/vendor:/this/project/path

go build -o bin/$(uname)_$(uname -m)/icd icd

Usage:
```
cf icd --create-connection <webhook> <cf app name>
```
Where <webhook> is the URL supplied shown in your toolchain configuration page for Jenkins broker
and <cf app name> is the Cloudfoundry application name which deployed successfully

```
cf icd --delete-connection <webhook> <cf app name>
```
Where <webhook> is the URL supplied shown in your toolchain configuration page for Jenkins broker
and <cf app name> is the Cloudfoundry application name which deployed removed.

## Build Instructions
The plugin can be built locally as a docker container, or by pushing to your forked branch and building with travis-ci.

### To Build Locally
```
git clone https://github.com/IBM/cf-icd-plugin /your/icd/dir
cd /your/icd/dir
./build

NOTE: ./build_bin contains your built binaries
```

### To Build with travis-ci
```
git clone https://github.com/IBM/cf-icd-plugin /your/icd/dir
cd /your/icd/dir
git remote add <your remote name> <your forked remote repo>
git tag -a <your release version id> -m 'my custom release version'
git push <your remote name> <your release version id>

NOTE:  Be sure to enable the remote repo in travis-ci.org after you link your github account with travis-ci.org account
```

# cf_icd_plugin
Sample CF plugin for POST of cloud foundry commands to registered webhook broker

export GOPATH=/this/project/path/vendor:/this/project/path

go build -o bin/$(uname)_$(uname -m)/icd icd

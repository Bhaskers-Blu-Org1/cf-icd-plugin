# cf_icd_plugin
Sample CF plugin for POST of cloud foundry commands to registered webhook broker

export GOPATH=<THIS DIR>/vendor:<THIS DIR>
go build -o bin/$(uname)_$(uname -m)/icd icd

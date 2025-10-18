SHELL := /bin/bash

gen-proto:
	cd external/plugin/v1/api/proto && buf lint && buf generate

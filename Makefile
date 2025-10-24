SHELL := /bin/bash

gen-proto: gen-proto-plugin gen-proto-api

gen-proto-api:
	cd external/plugin/v1/api/proto && buf lint && buf generate

gen-proto-plugin:
	cd api/proto && buf lint && buf generate

run-local-external_py_plugin: gen-proto
	uv run python cmd/external_py_plugin/main.py

run-local-external_go_plugin: gen-proto
	go run cmd/external_go_plugin/main.go

call-local-allgreetings:
	cd api/proto && grpcurl \
		-protoset <(buf build -o -) -plaintext \
		-d '{"name": "Moe"}' \
		localhost:8087 admin.v1.AdminService/AllGreetings | jq

call-local-allgreetingstreams:
	cd api/proto && grpcurl \
		-protoset <(buf build -o -) -plaintext \
		-d '{"name": "Moe"}' \
		localhost:8087 admin.v1.AdminService/AllGreetingStreams | jq

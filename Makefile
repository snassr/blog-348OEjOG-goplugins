SHELL := /bin/bash

gen-proto: gen-proto-plugin gen-proto-api

gen-proto-api:
	cd external/plugin/v1/api/proto && buf lint && buf generate 

gen-proto-plugin:
	cd api/proto && buf lint && buf generate 

call-local-allgreetings:
	cd api/proto && grpcurl \
		-protoset <(buf build -o -) -plaintext \
		-d '{"name": "Moe"}' \
		localhost:8087 admin.v1.AdminService/AllGreetings | jq

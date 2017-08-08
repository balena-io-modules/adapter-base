#! /usr/bin/bash

# Install dependencies
glide install
go install ./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go install ./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

# Generate golang stubs
protoc \
	-I ./adapter \
	-I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I /usr/local/include \
	./adapter/adapter.proto \
	--go_out=plugins=grpc:adapter

# Generate js stubs
protoc \
	-I ./adapter \
	-I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I /usr/local/include \
	./adapter/adapter.proto \
	--js_out=import_style=commonjs,binary:adapter

# Generate http stubs
protoc \
	-I ./adapter \
	-I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I /usr/local/include \
	./adapter/adapter.proto \
	--grpc-gateway_out=logtostderr=true:adapter

# Compile application
go build

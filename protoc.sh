#!/bin/sh

# Generate golang stubs
protoc \
	-I ./protos \
	-I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I /usr/local/include \
	./protos/update.proto \
	--go_out=plugins=grpc:update

protoc \
	-I ./protos \
	-I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I /usr/local/include \
	./protos/scan.proto \
	--go_out=plugins=grpc:scan

# Generate js stubs
protoc \
	-I ./protos \
	-I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I /usr/local/include \
	./protos/update.proto \
	--js_out=import_style=commonjs,binary:update

protoc \
	-I ./protos \
	-I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I /usr/local/include \
	./protos/scan.proto \
	--js_out=import_style=commonjs,binary:scan

# Generate http stubs
protoc \
	-I ./protos \
	-I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I /usr/local/include \
	./protos/update.proto \
	--grpc-gateway_out=logtostderr=true:update


protoc \
	-I ./protos \
	-I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I /usr/local/include \
	./protos/scan.proto \
	--grpc-gateway_out=logtostderr=true:scan

#!/bin/sh

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
	./protos/update.proto \
	--grpc-gateway_out=logtostderr=true:update

protoc \
	-I ./protos \
	-I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I /usr/local/include \
	./protos/scan.proto \
	--go_out=plugins=grpc:scan

protoc \
	-I ./protos \
	-I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I /usr/local/include \
	./protos/scan.proto \
	--grpc-gateway_out=logtostderr=true:scan

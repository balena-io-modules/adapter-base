#! /usr/bin/bash

glide install
go install ./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go install ./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
./protoc.sh
go build

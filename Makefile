BUILD_DIR ?= ./build

build:
	go build -o $(BUILD_DIR)/rkvs ./cmd/rkvs

%.pb.go: %.proto
	protoc $< --go_out=plugins=grpc:.

protoc_all: protoc_service

protoc_service: ./pkg/proto/service.pb.go

.PHONY: build
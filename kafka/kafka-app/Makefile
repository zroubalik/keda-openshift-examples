IMAGE ?= quay.io/zroubalik/kafka-app

ARCH       ?=amd64
TARGET_OS  ?=linux

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: build
build: producer consumer

.PHONY: producer
producer:
	go build -o kafkaproducerapp cmd/kafkaproducerapp/main.go

.PHONY: consumer
consumer:
	go build -o kafkaconsumerapp cmd/kafkaconsumerapp/main.go

.PHONY: image
image: image-build image-push

.PHONY: image-build
image-build:
	docker build -t ${IMAGE} .

.PHONY: image-push
image-push:
	docker push ${IMAGE}

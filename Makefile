
GOPATH:=$(shell go env GOPATH)
MODIFY=proto/

.PHONY: proto
proto:
    
	protoc --micro_out=${MODIFY} --go_out=${MODIFY} ${MODIFY}/cart-api.proto
    

.PHONY: build
build: proto

	go build -o cart-api-api *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t cart-api-api:latest

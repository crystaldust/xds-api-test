SHELL := /bin/bash
binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 vgo build -a -installsuffix cgo -ldflags '-w' -o xds-api-test

docker: binary
	docker build -t go-chassis/xds-api-test:v1 ./

all: docker binary
	./distribute-image.sh go-chassis/xds-api-test:v1
	kubectl apply -f ./xds-api-test.yaml

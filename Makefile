SHELL := /bin/bash
binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w' -o xds-api-test

docker: binary
	docker build -t juzhen/xds-api-test:v1 ./

all: docker binary
	distribute-image.sh juzhen/xds-api-test:v1
	istioctl kube-inject -f ./xds-api-test.yaml > xds-api-test.injected.yaml
	kubectl apply -f ./xds-api-test.injected.yaml


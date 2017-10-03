IMAGE_NAME=aws_swarm_labeler
VERSION=0.1
REGISTRY?=hub.docker.com

.PHONY: run

all: build push

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o aws_swarm_labeler
	docker build -t ${IMAGE_NAME} .
run:
	docker run -it -v /var/run/docker.sock:/var/run/docker.sock -v ${HOME}/.aws/credentials:/root/.aws/credentials ${IMAGE_NAME} /aws_swarm_labeler $(ARGS)
push:
	docker tag ${IMAGE_NAME} ${REGISTRY}/${IMAGE_NAME}:${VERSION} && docker push ${REGISTRY}/${IMAGE_NAME}:${VERSION}
	docker tag ${IMAGE_NAME} ${REGISTRY}/${IMAGE_NAME}:`git rev-parse --short HEAD` && docker push ${REGISTRY}/${IMAGE_NAME}:`git rev-parse --short HEAD`


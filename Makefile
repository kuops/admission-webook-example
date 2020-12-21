.PHONY: all

IMAGE_REPO ?= registry.cn-beijing.aliyuncs.com/kuops
IMAGE_NAME ?= adminssion-webhook
IMAGE_TAG ?= $(shell date +v%Y%m%d)

image: build-image push-image

build:
	@GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -installsuffix cgo -o build/_output/bin/$(IMAGE_NAME) .

build-image: build
	@echo "Building the docker image: $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)..."
	@docker build -t $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG) -f Dockerfile .

push-image: build-image
	@echo "Pushing the docker image for $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG) and $(IMAGE_REPO)/$(IMAGE_NAME):latest..."
	@docker tag $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_REPO)/$(IMAGE_NAME):latest
	@docker push $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)
	@docker push $(IMAGE_REPO)/$(IMAGE_NAME):latest

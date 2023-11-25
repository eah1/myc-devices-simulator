VERSION := 1.0
API_IMAGE_NAME := myc-device-simulator-api

myc-device-simulator-api:
	DOCKER_BUILDKIT=1 docker build \
		-f deployment/docker/myc-device-simulator-api.dockerfile \
		-t $(API_IMAGE_NAME):$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		--build-arg GOARCH=$(ARCH) \
		.
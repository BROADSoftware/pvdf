
# Image URL to use all building/pushing image targets
VERSION ?= latest
IMG ?= pvdf/pvscanner:$(VERSION)

# Currently, only unix executable are generated.
GOARGS = GOOS=linux GOARCH=amd64

default:
	@echo "Provide a target (build, img or push)"

# Build an executable for local test (Use only on linux system)
.PHONY: pvscanner
build:
	$(GOARGS) go build -o ./pvscanner

img:
	cd ..; docker build -f pvscanner/Dockerfile . -t ${IMG}

push: img
	docker push ${IMG}



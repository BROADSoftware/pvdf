
# Image URL to use all building/pushing image targets
VERSION ?= latest
IMG ?= pvdf/pvscanner:$(VERSION)

# Currently, only unix executable are generated.
GOARGS = GOOS=linux GOARCH=amd64

default:
	@echo "Provide a target (build, push or pvscanner)"

# Build an executable for local test (Use only on linux system)
pvscanner:
	$(GOARGS) go build -o ./pvscanner

build:
	docker build . -t ${IMG}

push: build
	docker push ${IMG}



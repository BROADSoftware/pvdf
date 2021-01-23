
# Image URL to use all building/pushing image targets
VERSION ?= latest
IMG ?= pvdf/pvdf:$(VERSION)

# Currently, only unix executable are generated.
GOARGS = GOOS=linux GOARCH=amd64

# Build an executable for local test (Use only on linux system)
pvdf:
	$(GOARGS) go build -o ./pvdf

build:
	docker build . -t ${IMG}

push: build
	docker push ${IMG}



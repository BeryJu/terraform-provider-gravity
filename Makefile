.SHELLFLAGS += -x -e

HOSTNAME=registry.terraform.io
NAMESPACE=BeryJu
NAME=gravity
BINARY=terraform-provider-${NAME}
VERSION=99999
OS_ARCH=darwin_amd64
TESTARGS=-v -p 1 -race -coverprofile=coverage.txt -covermode=atomic

default: gen

# Run acceptance tests
.PHONY: test
test:
	TF_ACC=1 go test $(TESTARGS) ./...
	go tool cover -html coverage.txt -o coverage.html

build: gen-api
	go build -o ${BINARY}

# The provider's schema uses EnvDefaultFunc, so tfplugindocs classifies any
# attribute whose variable happens to be set as optional rather than required.
# Clear all of them so generation does not depend on the ambient environment.
gen:
	golangci-lint run -v
	GRAVITY_URL="" GRAVITY_TOKEN="" GRAVITY_INSECURE="" go generate

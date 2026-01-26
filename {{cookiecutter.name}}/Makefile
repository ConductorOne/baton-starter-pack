GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
BUILD_DIR = dist/${GOOS}_${GOARCH}
GENERATED_CONF := pkg/config/conf.gen.go

ifeq ($(GOOS),windows)
OUTPUT_PATH = ${BUILD_DIR}/{{ cookiecutter.name }}.exe
else
OUTPUT_PATH = ${BUILD_DIR}/{{ cookiecutter.name }}
endif

# Set the build tag conditionally based on BATON_LAMBDA_SUPPORT
ifdef BATON_LAMBDA_SUPPORT
	BUILD_TAGS=-tags baton_lambda_support
else
	BUILD_TAGS=
endif

.PHONY: build
build: $(GENERATED_CONF)
	go build ${BUILD_TAGS} -o ${OUTPUT_PATH} ./cmd/{{ cookiecutter.name }}

$(GENERATED_CONF): pkg/config/config.go go.mod
	@echo "Generating $(GENERATED_CONF)..."
	go generate ./pkg/config

.PHONY: generate
generate: $(GENERATED_CONF)

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: update-deps
update-deps:
	go get -d -u ./...
	go mod tidy -v
	go mod vendor

.PHONY: add-dep
add-dep:
	go mod tidy -v
	go mod vendor

.PHONY: clean
clean:
	rm -rf dist/
	rm -f $(GENERATED_CONF)

.DEFAULT_GOAL := build

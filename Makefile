COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
USER=$(shell whoami)
TIMESTAMP_RFC3339 := $(shell date +%Y-%m-%dT%T%z)
BINARY=charite-data-donation

# Go Variables
VERSION=0.0.1
CILINT_VERSION := v1.22
PKG=github.com/d4l-data4life/charite-data-donation/pkg/config
LDFLAGS="-X '$(PKG).Version=${VERSION}' -X '$(PKG).Commit=${COMMIT}' -X '$(PKG).Branch=${BRANCH}' -X '$(PKG).BuildUser=${USER}'"
GOCMD=go
GOBUILD=$(GOCMD) build -ldflags ${LDFLAGS}
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
DOCKER_IMAGE=${BINARY}
CONTAINER_NAME=${BINARY}
PORT=4444

DB_IMAGE=postgres
DB_CONTAINER_NAME=${BINARY}-postgres
DB_PORT=5444

SRC = cmd/api/*.go

.PHONY: all help test test-race test-cover test-silent unit-test unit-test-docker vendor build run clean docker-build docker-database docker-run docker-stop kube-deploy kube-local db dr

## ----------------------------------------------------------------------
## Help: Makefile for app: charite-data-donation
## ----------------------------------------------------------------------

help:               ## Show this help (default)
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: lint unit-test-docker  ## Run tests from Jenkins - linting and unit

unit-test: test-cover        ## Run native unit tests

# Testing for race conditions needs much higher timeouts
test-race:          ## Run tests natively (-race)
	$(GOTEST) -timeout 120s -race ./...

test-cover:         ## Run tests natively (-cover) in verbose mode
	$(GOTEST) -timeout 15s -cover -covermode=atomic -v ./...

test-silent:         ## Run tests natively in non-verbose mode
	$(GOTEST) -timeout 15s ./...

vendor:             ## Download and tidy go depenencies
	@go mod tidy

build:              ## Build app
	$(GOBUILD) -o $(BINARY) $(SRC)

run:                ## Run app natively
	$(GOCMD) run $(SRC)

clean:              ## Remove compiled binary
	rm -f $(BINARY)
	rm -rf deploy/manifests
	-docker rm -f $(CONTAINER_NAME)
	-docker rm -f $(DB_CONTAINER_NAME)

lint:              ## Run golint
	docker build \
		--build-arg CILINT_VERSION=${CILINT_VERSION} \
		--build-arg GITHUB_USER_TOKEN \
		-t "$(DOCKER_IMAGE):lint" \
		-f build/lint.Dockerfile \
		.
	docker run --rm "$(DOCKER_IMAGE):lint"

docker-build db:    ## Build Docker image
	docker build \
		--build-arg GITHUB_USER_TOKEN \
		--build-arg APP_VERSION="$(VERSION)" \
		--build-arg GIT_COMMIT="$(COMMIT)" \
		--build-arg BUILD_DATE="$(TIMESTAMP_RFC3339)" \
		-t $(DOCKER_IMAGE) \
		-f build/Dockerfile \
		.

unit-test-docker:        ## Runs `make unit-test-native` inside the Docker
	docker build \
		--build-arg APP_VERSION="$(VERSION)" \
		--build-arg GIT_COMMIT="$(COMMIT)" \
		--build-arg BUILD_DATE="$(TIMESTAMP_RFC3339)" \
		--build-arg GITHUB_USER_TOKEN \
		-t $(DOCKER_IMAGE):test \
		-f build/test.Dockerfile \
		.
	docker run --rm $(DOCKER_IMAGE):test

docker-database:    ## Run database in Docker
	docker run --name $(DB_CONTAINER_NAME) -d \
		-e POSTGRES_DB=charite-data-donation \
		-e POSTGRES_PASSWORD=postgres \
		-p $(DB_PORT):5432 $(DB_IMAGE)

docker-run dr:      ## Run app in Docker. Configure connection to a DB using CHARITE_DATA_DONATION_DB_HOST and CHARITE_DATA_DONATION_B_PORT
	-docker run --name $(DB_CONTAINER_NAME) --rm -d \
		-e POSTGRES_DB=charite-data-donation \
		-e POSTGRES_PASSWORD=postgres \
		-p $(DB_PORT):$(DB_PORT) $(DB_IMAGE)
	docker run --name $(CONTAINER_NAME) --rm -t -d \
		-e CHARITE_DATA_DONATION_DB_HOST=host.docker.internal \
		-p $(PORT):$(PORT) \
		$(DOCKER_IMAGE)

docker-stop:        ## Stop Docker container with the app
	docker stop $(CONTAINER_NAME)
	docker stop $(DB_CONTAINER_NAME)

kube-deploy:        ## Deploy Docker container to local Kubernetes (check kubectl context beforehand)
	kubectl config use-context docker-desktop
	mkdir -p deploy/manifests
	@for file in deploy/templates/*.yaml; do\
		kubetpl render $$file \
			-i deploy/config/local.yaml \
			-o $${file/templates/manifests};\
	done
	@for file in deploy/local/*.yaml; do\
		kubetpl render $$file \
			--allow-fs-access \
			-i deploy/config/local.yaml \
			-o $${file/local/manifests};\
	done

	kubectl apply -f deploy/manifests/

kube-delete:        ## Delete local Kubernetes deployment (requires manifest files to exist)
	kubectl config use-context docker-desktop
	kubectl delete -f deploy/manifests/

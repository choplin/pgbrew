src_dir := /go/src/github.com/choplin/pgenv
go_srcs := $(shell find $(CURDIR) -name '*.go')

.PHONY: all clean test docker-test build-docker-image dev-test

all: pgenv

pgenv: $(go_srcs)
	go build -v

clean:
	rm -f pgenv docker/pgenv

test:
	docker run --rm -v $(CURDIR):$(src_dir) -w $(src_dir) choplin/pgenv-test-env make docker-test

docker-test:
	go get -d ./...
	go build -v
	go test ./...

build-docker-image: docker/pgenv
	docker build --rm -t choplin/pgenv-test-env docker

docker/pgenv: $(go_srcs)
	GOOS=linux GOARCH=amd64 go build -o docker/pgenv

dev-test:
	docker run --rm -v $(GOPATH)/src:/go/src -w $(src_dir) choplin/pgenv-test-env go test --short ./...

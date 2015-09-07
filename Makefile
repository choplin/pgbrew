src_dir := /go/src/github.com/choplin/pgenv

.PHONY: test docker-test build-docker-image dev-test

test:
	docker run --rm -v $(CURDIR):$(src_dir) -w $(src_dir) choplin/pgenv-test-env make docker-test

docker-test:
	go get -d ./...
	go build -v
	go test ./...

build-docker-image:
	docker build --rm -t choplin/pgenv-test-env docker

dev-test:
	docker run --rm -v $(GOPATH)/src:/go/src -w $(src_dir) choplin/pgenv-test-env go test --short ./...

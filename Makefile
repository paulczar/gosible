PKGS=$(go list ./... | grep -v /vendor)
CI_PKGS=$(go list ./... | grep -v /vendor | grep -v test)
FMT_PKGS=$(go list -f {{.Dir}} ./... | grep -v vendor | grep -v test | tail -n +2)
GIT_SHA=$(git rev-parse --verify HEAD)
VERSION=$(cat VERSION)
PWD=$(pwd)

default: compile

compile: ## Create the gosible executable in the ./bin directory.
	go build -o bin/gosible main.go

install: ## Create the kubicorn executable in $GOPATH/bin directory.
	install -m 0755 bin/gosible ${GOPATH}/bin/gosible

test: test test-vagrant test-ping test-adhoc

test-vagrant:
	go run main.go playbook run -e tests/functional/environment \
	  --provisioner=vagrant \
    tests/functional/playbook/ping.yml  --become

test-ping:
	go run main.go ping -e tests/functional/environment

test-adhoc:
	go run main.go adhoc -e tests/functional/environment --command "hostname"

lint:
	golint $(PKGS)

gofmt:
	echo "Fixing format of go files..."; \
	for package in $(FMT_PKGS); \
	do \
		gofmt -w $$package ; \
		goimports -l -w $$package ; \
done

dep:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
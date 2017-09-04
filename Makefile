default: compile

compile: ## Create the gosible executable in the ./bin directory.
	go build -o bin/gosible main.go

test: test-vagrant test-ping test-adhoc

test-vagrant:
	go run main.go playbook run -e tests/functional/environment \
	  --provisioner=vagrant \
    tests/functional/playbook/ping.yml  --become

test-ping:
	go run main.go ping -e tests/functional/environment

test-adhoc:
	go run main.go adhoc -e tests/functional/environment --command "hostname"		
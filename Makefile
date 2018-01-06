APP = reminder
VERSION = 
REVISION = 

.PHONY: install
install:
	go install ./...

.PHONY: build
build:
	go build -o bin/darwin/$(APP) ./cmd

.PHONY: test
test: test-units test-integration

.PHONY: test-units
test-units:
	go test -v ./...

.PHONY: test-integration
test-integration:
	@echo "integration tests should start"


.PHONY: start
start:
	./bin/darwin/$(APP) & echo $! > $(APP).pid

.PHONY: stop
stop:
	cat $(APP).pid | kill -9

.PHONY: clean
clean:
	rm bin/*
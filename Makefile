APP = reminder
VERSION = 
REVISION = 

.PHONY: install
install:
	go install ./...

.PHONY: build
build:
	go build -o bin/darwin/$(APP) ./cmd/$(APP)

.PHONY: test
test: test-units test-integration

.PHONY: test-units
test-units:
	go test -v ./...

.PHONY: test-integration
test-integration:
	@echo "integration tests should start"

.PHONY: clean
clean:
	rm bin/*
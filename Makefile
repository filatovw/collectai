APP = reminder
PLATFORM ?= darwin
VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT)"

.PHONY: install
install:
	go install ./...

.PHONY: build
build:
	go build $(LDFLAGS) -o bin/$(PLATFORM)/$(APP) ./cmd/$(APP)

.PHONY: test
test:
	go test -v ./... -race

.PHONY: clean
clean:
	rm -rf bin/*
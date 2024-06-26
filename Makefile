REPOSITORY := github.com/np-guard/cloud-resource-collector
EXE:=collector

mod: go.mod
	@echo -- $@ --
	go mod tidy
	go mod download

fmt:
	@echo -- $@ --
	goimports -local $(REPOSITORY) -w .

lint:
	@echo -- $@ --
	CGO_ENABLED=0 go vet ./...
	golangci-lint run

precommit: mod fmt lint

build:
	@echo -- $@ --
	CGO_ENABLED=0 go build -o ./bin/$(EXE) ./cmd/collect

test:
	@echo -- $@ --
	go test ./... -v -cover -coverprofile collect.coverprofile
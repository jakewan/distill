.DEFAULT_GOAL := local-dev-all

.PHONY: go-staticcheck
go-staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck -go 1.21 ./...

.PHONY: go-test
go-test:
	go test ./...

.PHONY: go-vet
go-vet:
	go vet ./...

.PHONY: local-dev-all
local-dev-all: go-test go-vet go-staticcheck

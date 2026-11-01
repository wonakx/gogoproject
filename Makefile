.PHONY: generate check-api test docs

generate:
	go generate ./api

check-api:
	sh scripts/check-openapi.sh

test: check-api
	go test -race ./...
	go vet ./...

docs:
	go run ./cmd/apidocs

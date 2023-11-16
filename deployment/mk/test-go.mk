
tidy:
	go mod tidy
	go mod vendor

lint:
	golangci-lint version
	golangci-lint linters
	golangci-lint run

test:
	go test -coverprofile=profile.cov ./... -p 2
	go tool cover -func profile.cov
	go vet ./...
	gofmt -l .
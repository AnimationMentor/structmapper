
test:
	go test -v ./...

lint:
	golangci-lint -v run
test:
	go test
	go vet
	golint ./...

.PHONY: test

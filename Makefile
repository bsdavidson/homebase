default: test build


build:
	-mkdir build
	env GOOS=darwin GOARCH=amd64 go build -o build/homebase_darwin_amd64 cmd/homebase/main.go
	env GOOS=linux GOARCH=amd64 go build -o build/homebase_linux_amd64 cmd/homebase/main.go
	env GOOS=linux GOARCH=386 go build -o build/homebase_linux_386 cmd/homebase/main.go
	env GOOS=windows GOARCH=amd64 go build -o build/homebase_windows_amd64 cmd/homebase/main.go
	env GOOS=windows GOARCH=386 go build -o build/homebase_windows_386 cmd/homebase/main.go

test:
	go test
	go vet
	golint ./...

.PHONY: build test

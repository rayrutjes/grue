.PHONY: test
test:
	go test ./...

.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build -o .build/grue-linux-amd64

.PHONY: mac
mac:
	GOOS=darwin GOARCH=amd64 go build -o .build/grue-darwin-amd64

.PHONY: clean
clean:
	go clean
	rm -rf .build

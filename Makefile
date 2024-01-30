build:
	mkdir bin
	go mod init note
	go get -v -d ./...
	go mod tidy
	go build -o ./bin ./...

test:
	go test ./... | grep -v 'no test files'

lint:
	golangci-lint -c .golangci.yml run ./...

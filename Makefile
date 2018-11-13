BINARY_NAME=new-backend

all: test build
build:
	go build -o $(BINARY_NAME) -v
test:
	go test -v ./...
clean:
	go clean
	rm $(BINARY_NAME)


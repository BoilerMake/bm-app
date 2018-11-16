# Add any other binaries to build here, seperated by a space
TARGETS := serve

all: test build serve
	
build:
	for target in $(TARGETS); do \
		go build -o bin/$$target ./cmd/$$target; \
	done

test:
	# Runs every tests (*_test.go)
	go test -v ./...

clean:
	go clean
	rm -rf ./bin

serve:
	./bin/serve

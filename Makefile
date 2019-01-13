# TODO this should be replaced with make.go which you can run with `go run make.go`
# I know we've got some windows bois (🤮 it's 2018, use unix 😡), and this should solve their woes
# ^ That doesn't mean we shouldn't have a Makefile though! running `make` is still faster than `go run make.go`

# Add any other binaries to build here, seperated by a space
TARGETS := server

INFO_STR=[INFO]

all: test build server
build:
	@for target in $(TARGETS); do \
		echo $(INFO_STR) building binary \"$$target\"; \
		go build -o bin/$$target ./cmd/$$target; \
	done

test:
	@echo $(INFO_STR) running tests
	@go test -v ./...

clean:
	@echo $(INFO_STR) cleaning dependencies and removing binaries
	@go clean
	@go mod tidy
	@rm -rf ./bin

server:
	@echo $(INFO_STR) running bin/server
	@./bin/server

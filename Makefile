# Add any other binaries to build here, seperated by a space
TARGETS := server

INFO_STR=[INFO]

build:
	@for target in $(TARGETS); do \
		echo $(INFO_STR) building binary \"$$target\"; \
		go build -o bin/$$target ./cmd/$$target; \
	done

dev:
	@echo $(INFO_STR) starting dev environment...
	@docker-compose -f deploy/docker-compose.base.yml -f deploy/docker-compose.dev.yml up

dev-rebuild:
	@echo $(INFO_STR) rebuilding dev environment...
	@docker-compose -f deploy/docker-compose.base.yml -f deploy/docker-compose.dev.yml up --build --force-recreate

dev-cleanup:
	@echo $(INFO_STR) removing dev environment...
	@docker-compose -f deploy/docker-compose.base.yml -f deploy/docker-compose.dev.yml rm --stop

test:
	@echo $(INFO_STR) starting test environment...
	@docker-compose -f deploy/docker-compose.base.yml -f deploy/docker-compose.testing.yml up -d
	@docker-compose -f deploy/docker-compose.base.yml -f deploy/docker-compose.testing.yml exec backend go test -v /backend/...

clean:
	@echo $(INFO_STR) cleaning dependencies and removing binaries...
	@go clean
	@go mod tidy
	@rm -rf ./bin


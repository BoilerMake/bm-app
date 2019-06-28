# Add any other binaries to build here, seperated by a space
TARGETS := server

INFO_STR=[INFO]

dev:
	@echo $(INFO_STR) starting dev environment...
	@docker-compose -f deploy/docker-compose.default.yml -f deploy/docker-compose.dev.yml up

test:
	@echo $(INFO_STR) starting test environment...
	@# Because of how docker handles .env files (different than env_files‽), we
	@# need to run the docker-compose command inside the directory with our
	@# test .env file.
	@#
	@# Why is this one big command? Make runs separate commands in their own
	@# subshell (oof @252), so if they weren't together then the following command
	@# would no longer in the correct directory after cding.
	@cd deploy/test && \
		docker-compose -f ../docker-compose.default.yml -f docker-compose.test.yml up -d && \
		docker-compose -f ../docker-compose.default.yml -f docker-compose.test.yml exec bm-app go test -v /bm-app/...
	@# Also, see comment above for why we don't need to cd out

build:
	@for target in $(TARGETS); do \
		echo $(INFO_STR) building binary \"$$target\"; \
		go build -o bin/$$target ./cmd/$$target; \
	done

dev-force-rebuild:
	@echo $(INFO_STR) rebuilding dev environment...
	@docker-compose -f deploy/docker-compose.default.yml -f deploy/docker-compose.dev.yml up --build --force-recreate

dev-cleanup:
	@echo $(INFO_STR) removing dev environment...
	@docker-compose -f deploy/docker-compose.default.yml -f deploy/docker-compose.dev.yml rm --stop

review-app:
	@echo $(INFO_STR) starting dev environment...
	@dh -f deploy/docker-compose.default.yml -f deploy/docker-compose.dev.yml up

clean:
	@echo $(INFO_STR) cleaning dependencies and removing binaries...
	@go clean
	@go mod tidy
	@rm -rf ./bin

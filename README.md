# bm-app

## Development instructions
### Getting started
- [Install docker and docker-compose](https://docs.docker.com/compose/install/)
- [Install Node and npm](https://nodejs.org/en/)
- Copy and setup variables in `.env`
	- `$ cp .env.example .env`
  - Then make any env changes you need

### Running the server
- `$ make dev`
	- Sets up the server and automatically recompiles on file changes
- `$ make test`
	- Runs our test suite

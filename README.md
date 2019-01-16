# new-backend

## Development instructions
### Getting started
- [Install Go](https://golang.org/doc/install)
- Copy and setup variables in `.env`
	- `$ cp .env.example .env`
  - Then make any env changes you need
- Install [PostgreSQL](https://www.postgresql.org/)
	- If you installed with brew, run `$ brew services start postgresql`

### Bootstrapping the database
- Create the database
	- `$ createdb boilermake`
- Run migrations using [goose](https://github.com/pressly/goose)
	- Install goose
		- `$ go get -u github.com/pressly/goose/cmd/goose`
	- `$ cd migrations`
	- `$ goose postgres YOUR_CONNSTR up`
		- For example it could be `$ goose postgres "user=cj dbname=boilermake sslmode=disable" up`

### Running the server
- `$ make`
  - Which will test, build, and run the server    

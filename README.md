# Ugo

Backend services for Ujo implemented in Go. Currenty, the primary service is the Ujo API Gateway which
provides an entrypoint for interacting with an Ujo indexer node. The current implementation uses a
PostgreSQL instance behind the API server as the database.
Please see below for instructions to get started running the application
locally. There is additional information in regard to running the tests and the project structure.

### Setup environment

You will need GO > 1.11 in order to run Ugo.  

Ensure that the `GOPATH` is set in .profile or .bash_profile
More info [here](https://golang.org/doc/code.html#GOPATH).

```sh
export GOPATH=$HOME/<your development path to>/go
PATH=$PATH:$GOPATH/bin
PATH=$PATH:/usr/local/go/bin
```

Clone this repository somewhere.  You can either clone it inside or outside of your `GOPATH`.  If you clone it inside your `GOPATH`, you'll need to do something like the following:

```sh
$ mkdir -p $GOPATH/src/github.com/consensys
$ cd $GOPATH/src/github.com/consensys
$ git clone git@github.com:ConsenSys/ugo.git

# Because we're cloning inside of GOPATH, we have to tell Go that we still
# want it to use gomod for dependency management:
$ echo 'export GO111MODULE=on' >> ~/.profile
# You'll need to restart your terminal session after the above command ^

$ cd ugo
```

If you clone outside of `GOPATH`, life is a little easier:

```sh
$ cd ~/projects
$ git clone git@github.com:ConsenSys/ugo.git
$ cd ugo
```

### Install dependencies

```
go mod download
```

Add environment variables. Replace username with the one created in postgres.

```
export DATABASE_URL=postgres://<username>:@localhost/ujo?sslmode=disable  
export UJO_API_SECRET=secret
```

You can add these exports to your .profile or .bash_profile too for future use.

### Setup Database

Install PostgreSQL. This can be done with homebrew on MacOS

```
brew update
```

```
brew install postgres
```

Create the database

```
createdb ujo
```

The postgress database needs to be run on its own in a separate thread.

```
postgres
```

sql-migrate needs to be installed from inside the GOPATH in order to get the database up to speed. Following instructions from: https://github.com/rubenv/sql-migrate, install it as follows:

```
go get -v github.com/rubenv/sql-migrate/...
```
Then: Change into the repository root directory and run sql-migrate to run database migrations. s

```
sql-migrate up
```

### Run

To start the API service, run:

```
go run ./cmd/api/api.go
```

### Build

To build a binary of the API into the bin dir, run:

```
go build -o bin/api ./cmd/api
```

### Watch

For rapid development, the go-watcher utility is helpful for auto-building the source files whenever changes are made. Get the package via:

```
go get github.com/canthefason/go-watcher
```

Install the binary under the `go/bin` directory:

```
go install github.com/canthefason/go-watcher/cmd/watcher
```

Finally, run the watcher:

```
watcher -c config -run github.com/consensys/ugo/cmd/api -watch github.com/consensys/ugo
```

### Make

To build a binary of the application run make from the root directory of the project

```
Make build
```

## Testing

To set up for the tests
```
createuser ubuntu
```

```
createdb ujo-test
```

To run the tests

```
go test ./pkg/*/ -v
```

To run the tests in serial, run:
```
go test -p 1 ./pkg/*/ -v
```

## Project Structure

The structure of the codebase follows the best practices outlined by Peter Bourgon's:

[Go best practices, 6 years in](https://peter.bourgon.org/go-best-practices-2016/#repository-structure).

There are two top-level directories, pkg and cmd. Underneath pkg, there are directories for each of the project's libraries. Underneath cmd, there are directories for each of your binaries. All of the Go code should live exclusively in one of these locations. In the root directory there are various configuration / continuous integration files, and additional directories for database code and migrations. A simplified example of this looks like the following:


```
github.com/consensys/ugo/
  circle.yml
  Dockerfile
  cmd/
    foosrv/
      main.go
    foocli/
      main.go
  pkg/
    fs/
      fs.go
      fs_test.go
      mock.go
      mock_test.go
    merge/
      merge.go
      merge_test.go
    api/
      api.go
      api_test.go
```

## Fixtures

Load prepopulated data

```
go run ./db/loaddata.go
```

Drop all data

```
go run ./db/dropdata.go
```

**Note**:  These loaddata migrates the database up so it will fail if tables already exist. Drop the db with `dropdb <dbname>` and recreate it with `createdb <dbname>`.

## Dependencies

This project currently uses the [Go 1.11 Modules](https://github.com/golang/go/wiki/Modules)
tool to manage dependencies.

## Infrastructure

At the moment the code is hosted via Heroku
to ease devops requirements on our small engineering team. It is likely that in the future we will migrate
services to more flexible infrastructure.

# {{ cookiecutter.project_name }}

## Selected Tables

{% for table in cookiecutter.selected_tables["values"] %}
- {{ table }}
{% endfor %}

A small skeleton project for a RESTful API using Go with Echo that uses a PostgreSQL database under the hood

## Requirements

- [Git](https://git-scm.com/)
- [Go](https://go.dev/) >= 1.22

## Initial setup

```shell
git clone git@github.com:neowinx/postgres-go-echo-restful.git
cd postgres-go-echo-htmx-bulma
```

## How to use

Once you've cloned the project you can run it by either...

### ... compiling and running the generated binary

```shell
go build ./web/server.go
export DATABASE_URI="postgresql://postgres:postgres@localhost:5432/super"
./server
```

### ... or running it directly

```shell
export DATABASE_URI="postgresql://postgres:postgres@localhost:5432/super"
go run ./web/server.go
```

### ... or by using [air](https://github.com/air-verse/air)

```shell
export DATABASE_URI="postgresql://postgres:postgres@localhost:5432/super"
air
```

and then you can go to this url to access the swagger documentation `http://127.0.0.1:8080/swagger/index.html`


## SQLC generation

This project uses [sqlc](https://sqlc.dev/) for the database abstraction. In order to generate the queries just edit the `querys.sql` file and then execute: 

```shell
sqlc generate
```

to generate the interfaces and magic stuff in `pkg/db`

## Swagger Generation

The swagger documentation is generated using [swaggo](https://github.com/swaggo/swag) and [swagger-echo](https://github.com/swaggo/echo-swagger) respectively. This works by
issuing the next command in order to generate the documentation based on the comments in the go files:

```bash
swag init -g server.go -d ./web,./internal
```

NOTE: This needs to be executed everytime you need to update the swagger documentation

## Docker

There is a `compose.yml` configuration to use with [docker-compose](https://docs.docker.com/compose/)

Just run this command and it will build and run database and application images:

```shell
docker compose up -d
```

and then go to this url in your browser http://127.0.0.1:8080/swagger/index.html

once you quit playing around with it you can issue the next command to stop the services and delete the images

```shell
docker compose down -v --rmi local
```

## Enjoy

Pedro

# REST API template 【Golang, Docker】
This Repository is RESTAPI quick start template.

☑️ docker-compose

☑️ mysql

☑️ [golang-migrate](https://github.com/golang-migrate/migrate)

☑️ [sqlboiler](https://github.com/volatiletech/sqlboiler)

☑️ [wire](https://github.com/google/wire)

☑️ handler template

☑️ Makefile

## Development
### Setup
you need to install docker, docker-compose, golang
```sh
$ git clone (URL)
$ make setup
```

### Run
```sh
$ make start
```

### Build
```sh
# output bin/api-server
$ make build
```

## Generate
```sh
$ make gen
```

### Migrate
```sh
# up
$ make migrate-up
# down
$ make migrate-down
```

# 時間管理 slack app サーバAPI
This Repository is quick start template for REST API.

☑️ docker-compose

☑️ mysql

☑️ [golang-migrate](https://github.com/golang-migrate/migrate)

☑️ [sqlboiler](https://github.com/volatiletech/sqlboiler)

☑️ [wire](https://github.com/google/wire)

☑️ [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)

☑️ handler template

☑️ Makefile

☑️ GitHub Action

## Development
### Setup
- you need to install docker, docker-compose, golang.
- you need to setting these path.
```sh
$ git clone (URL)
$ make setup
```

### Run
```sh
$ docker compose start
$ go run ./cmd/server.go
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

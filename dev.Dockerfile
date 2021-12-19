FROM ubuntu:18.04

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y curl gnupg2 vim
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.11.0/migrate.linux-amd64.tar.gz | tar xvz
RUN mv ./migrate.linux-amd64 /usr/bin/migrate

RUN mkdir app
WORKDIR app 
ADD . .

RUN chmod 777 ./script/migrate-create
RUN chmod 777 ./script/migrate-up
RUN chmod 777 ./script/migrate-down

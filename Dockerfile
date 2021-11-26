FROM golang:latest

# RUN apt-get update && apt-get upgrade -y
# RUN apt-get install -y curl gnupg2 make
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.11.0/migrate.linux-amd64.tar.gz | tar xvz
# RUN mv ./migrate.linux-amd64 /usr/bin/migrate
# 
# RUN mkdir app
# WORKDIR app 
# ADD . .
# 
# ENV DSN OOOOOOOOOOOOOOOOOO
# ENV PATH $(go env GOPATH)/bin:$PATH
# 
# RUN make setup
# RUN make migrate
# 
# ENV PATH /app/bin:$PATH

##  App
##-----------------------------------------------b
FROM golang:latest as build-env

ENV APP_ROOT $GOPATH/src/time-management-go
RUN ln -s $APP_ROOT/ /app
WORKDIR /app 
COPY . $APP_ROOT/
RUN GOOS=linux	go build -o ./bin ./cmd/server.go

##  Runtime build stage
##-----------------------------------------------
FROM debian:10.8-slim
RUN mkdir app

COPY --from=build-env /app/bin /app/bin
ENV PATH /app/bin:$PATH
RUN chmod a+x bin/*

EXPOSE 8000
CMD ["server"]

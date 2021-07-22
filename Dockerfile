FROM golang:1.17-rc-alpine3.13 AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY ./src/go.mod ./
COPY ./src/*.go ./

RUN GOOS=linux GOARCH=arm GOARM=7 go mod tidy
RUN GOOS=linux GOARCH=arm GOARM=7 go build -o ./slackit .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN touch CREATED.txt

WORKDIR /
COPY --from=build ./src/slackit ./
EXPOSE 8080
CMD ["./slackit"]


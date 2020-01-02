FROM golang:1.13-alpine AS build

WORKDIR $GOPATH/src/github.com/brennerm/azure-keyvault-external-secret
COPY go.mod .
COPY go.sum .
COPY Makefile .
RUN go mod download 
RUN go mod vendor
COPY controller.go .
COPY main.go .
COPY crd crd
COPY keyvault keyvault

RUN apk update && apk add make
RUN make

FROM alpine:latest
COPY --from=build /go/src/github.com/brennerm/azure-keyvault-external-secret/azure-keyvault-external-secret .
ENTRYPOINT ./azure-keyvault-external-secret

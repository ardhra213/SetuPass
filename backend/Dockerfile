# syntax=docker/dockerfile:1
ARG BASE_IMAGE=mosipdev/esignet:develop-go
FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY thunderid ./thunderid
COPY esignet-service ./esignet-service
WORKDIR /src/esignet-service
RUN go mod edit -replace github.com/thunder-id/thunderid=/src/thunderid && go mod download
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/esignet.exe ./cmd/esignet
FROM ${BASE_IMAGE}
COPY --from=build /out/esignet.exe /home/mosip/esignet.exe
COPY esignet-service/data /home/mosip/data

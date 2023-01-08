FROM golang:1-alpine as build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY cmd cmd
COPY internal internal

RUN CGO_ENABLED=0 go build -o /ethereum-healthmon cmd/healthmon/main.go

FROM alpine

COPY --from=build /ethereum-healthmon /usr/bin/

EXPOSE 21171

ENTRYPOINT ["ethereum-healthmon"]

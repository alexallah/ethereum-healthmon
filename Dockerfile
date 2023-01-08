FROM golang:1-alpine as build

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 go build -o /ethereum-healthmon cmd/healthmon/main.go

FROM alpine

COPY --from=build /ethereum-healthmon /usr/bin/

EXPOSE 21171

ENTRYPOINT ["ethereum-healthmon"]

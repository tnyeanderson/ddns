# build container
FROM golang:1.22 as build

WORKDIR /src

COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -o /opt/ddns .

# main container
FROM alpine

COPY --from=build /opt/ddns /usr/local/bin/ddns

ENTRYPOINT ["ddns"]

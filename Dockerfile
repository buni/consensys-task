FROM golang:1.17-bullseye as build

WORKDIR /build
ADD . /build
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -a -v -o service ./cmd/api/
RUN curl -o root.crt -O https://cockroachlabs.cloud/clusters/fa5249b5-e2b3-4e43-a224-765f2ef2c439/cert

FROM alpine:3.15.0
COPY --from=build /build/service /  
ENTRYPOINT ["/service"]


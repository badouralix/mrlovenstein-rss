##
## Build
##

FROM golang:1.17-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o mrlovenstein-rss

##
## Deploy
##

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /app/mrlovenstein-rss /mrlovenstein-rss

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/mrlovenstein-rss"]

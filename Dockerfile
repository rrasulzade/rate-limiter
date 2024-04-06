ARG GOLANG_VERSION=1.17
ARG ALPINE_VERSION=3.19

## Build stage
FROM golang:${GOLANG_VERSION}-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app

# final stage
FROM alpine:${ALPINE_VERSION} AS final

WORKDIR /app

COPY --from=build /app/app .
COPY --from=build /app/config/config.json ./config/config.json

CMD ["./app"]

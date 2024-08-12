FROM golang:1.22-bullseye as builder

WORKDIR /build

COPY go.sum go.mod ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./sonar-to-confluence ./cmd

FROM {DOCKER_HOST}/standard-images/ubuntu:focal

WORKDIR /src

COPY --from=builder  /build/sonar-to-confluence .
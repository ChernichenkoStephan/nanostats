# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY . ./

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./bin/nanostats ./cmd/nanostats

CMD [ "./bin/nanostats", "-config", "config.yaml" ]
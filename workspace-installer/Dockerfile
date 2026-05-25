FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /vpsik .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates restic docker-cli
COPY --from=builder /vpsik /usr/local/bin/vpsik
ENTRYPOINT ["vpsik"]

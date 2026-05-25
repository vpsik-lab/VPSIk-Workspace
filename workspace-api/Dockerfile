FROM golang:1.22-alpine AS builder
WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum ./
RUN GOPROXY=https://proxy.golang.org,direct go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /workspace-api .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /workspace-api /usr/local/bin/workspace-api
EXPOSE 8081
ENTRYPOINT ["workspace-api"]

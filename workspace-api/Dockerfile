FROM golang:1.22-bookworm AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /workspace-api .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /workspace-api /usr/local/bin/workspace-api
EXPOSE 8081
ENTRYPOINT ["workspace-api"]

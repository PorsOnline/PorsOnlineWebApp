FROM golang:1.23-alpine AS builder
WORKDIR /app

# First, copy only the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies in a separate layer
RUN GO111MODULE=on GOPROXY=https://goproxy.cn,direct go mod download

# Now copy the rest of the project
COPY . .

# Build the application
RUN GO111MODULE=on GOPROXY=https://goproxy.cn,direct go build -o server ./cmd/votingSystem/main.go

FROM alpine:latest
RUN apk --no-cache add tzdata
ENV TZ=Asia/Tehran
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /app/server .
COPY --from=builder /app/certs .
EXPOSE 8080
CMD ["./server"]
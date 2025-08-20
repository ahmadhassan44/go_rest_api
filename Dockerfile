FROM golang:1.24.5 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM alpine:latest
COPY --from=builder /app/bin/go_rest_api /app/bin/go_rest_api
WORKDIR /app
EXPOSE 3000
CMD [ "./bin/go_rest_api" ]
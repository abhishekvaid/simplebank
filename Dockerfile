# Build Stage 
FROM golang:1.22.4-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o simplebank main.go


# RUN Stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/simplebank /app/
COPY --from=builder /app/app.env /app/
CMD ["/app/simplebank"]

EXPOSE 8081
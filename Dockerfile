# Build Stage 
FROM golang:1.22.4-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-arm64.tar.gz | tar xvz
RUN go build -o simplebank main.go


# RUN Stage
FROM alpine
WORKDIR /app

# Copy the built binary and other necessary files
COPY --from=builder /app/simplebank .
COPY --from=builder /app/app.env .
COPY --from=builder /app/startup.sh .
COPY --from=builder /app/migrate /usr/bin
COPY --from=builder /app/db/migration ./db/migration

CMD ["/app/simplebank"]
ENTRYPOINT [ "/app/startup.sh" ]

EXPOSE 8081
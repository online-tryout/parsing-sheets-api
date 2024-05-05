FROM golang:1.22.2-alpine3.19 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates curl
WORKDIR /app
RUN mkdir -p /app
RUN curl -o /app/sheets-key.json -L https://storage.googleapis.com/online-tryout/sheets-key.json
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .
COPY --from=builder /app/Makefile .
COPY --from=builder /app/db ./db

CMD ["./main"]
FROM golang:1.22.0-alpine

WORKDIR /app

COPY . .

# install psql
RUN apk update && \
    apk add --no-cache postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

RUN go build -o app ./cmd/banner/main.go

EXPOSE 8080
EXPOSE 8082

CMD ["./app"]
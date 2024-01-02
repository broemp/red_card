FROM golang:alpine AS BUILD

WORKDIR /build

COPY go.mod .
RUN go mod download

COPY . .

RUN go build -tags=viper_bind_struct -o /redCard cmd/main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

FROM busybox
WORKDIR /app
COPY db/migration ./migration
COPY start.sh ./start.sh
RUN chmod +x ./start.sh
COPY --from=BUILD /build/migrate /app/migrate
COPY --from=BUILD /redCard .

EXPOSE 3000
CMD ["/app/redCard"]
ENTRYPOINT [ "/app/start.sh" ]

# ---------------------------
# Stage 1: Build
# ---------------------------
FROM golang:1.25.1-alpine AS builder

WORKDIR /srv

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /srv/user/cmd
RUN go build -o /srv/user/main .

# ---------------------------
# Stage 2: Run
# ---------------------------
FROM alpine:3.20

WORKDIR /srv

COPY --from=builder /srv/user/main ./main
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /srv/user/migrations /srv/migrations

COPY --from=builder /srv/certs/private.pem /srv/certs/private.pem
COPY --from=builder /srv/certs/public.pem /srv/certs/public.pem

ENV PATH="/usr/local/bin:${PATH}"

EXPOSE 80

CMD ["./main"]

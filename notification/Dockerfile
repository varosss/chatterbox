# ---------------------------
# Stage 1: Build
# ---------------------------
FROM golang:1.25.1-alpine AS builder

WORKDIR /srv

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /srv/notification/cmd
RUN go build -o /srv/notification/main .

# ---------------------------
# Stage 2: Run
# ---------------------------
FROM alpine:3.20

WORKDIR /srv

COPY --from=builder /srv/notification/main ./main

COPY --from=builder /srv/certs/public.pem /srv/certs/public.pem

ENV PATH="/usr/local/bin:${PATH}"

EXPOSE 80

CMD ["./main"]

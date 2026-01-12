FROM golang:1.24-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG MAIN_PKG=./cmd/app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/beanefits ${MAIN_PKG}

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/beanefits /app/beanefits
COPY --from=builder /src/db/migrations /app/migrations

EXPOSE 8080
ENTRYPOINT ["/app/beanefits"]

FROM golang:1.24 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build \
    --ldflags '-extldflags "-static"' \
    -o login-api \
    cmd/main.go

FROM scratch

WORKDIR /

COPY --from=build /app/login-api /bin/login-api

USER 1001:1001
CMD ["/bin/login-api"]
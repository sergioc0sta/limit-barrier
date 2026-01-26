FROM golang:1.24-alpine AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/limiter-barrier ./cmd/server/main.go

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=build /out/limiter-barrier /app/limiter-barrier
COPY configs /app/configs

ENTRYPOINT ["/app/limiter-barrier"]

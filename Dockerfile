FROM golang:1.15-alpine AS build
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o netrics ./cmd/server/main.go

FROM alpine:3 AS final
WORKDIR /app
COPY --from=build /build/netrics .
EXPOSE 9091
ENTRYPOINT ["./netrics"]
CMD ["-addr", ":9091", "-endpoint", "/metrics", "-sc", "3", "-interval", "30m"]
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/server /app/server
RUN adduser -D -H -u 65532 nonroot
USER nonroot
EXPOSE {{HTTP_PORT}}
ENV HTTP_PORT={{HTTP_PORT}}
CMD ["/app/server"]

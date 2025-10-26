# Build stage
FROM golang:1.25-alpine AS build
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./ 
RUN go env -w GOPROXY=https://proxy.golang.org
COPY . .
RUN go build -o /packcalc ./cmd/packcalc

# Final stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /packcalc /usr/local/bin/packcalc
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/usr/local/bin/packcalc"]
FROM golang:1.24-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH

RUN go build -trimpath -ldflags="-s -w" -o /out/app .

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/app /app/app
COPY --from=builder /src/assets /app/assets

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/app/app"]

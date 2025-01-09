# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build
COPY . .

RUN go mod download && \
  go vet -v ./... && \
  go test -v ./... && \
  GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 \
  go build -o ./ddnsgo ./cmd/app/main.go

FROM gcr.io/distroless/base-nossl-debian12

COPY --from=builder /build/ddnsgo /

CMD ["/ddnsgo"]

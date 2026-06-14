FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /feedbackboard ./cmd

FROM debian:12-slim
RUN groupadd -r nonroot && useradd -r -g nonroot nonroot
WORKDIR /
COPY --from=builder /feedbackboard /feedbackboard
COPY --from=builder /app/ui /ui
COPY --from=builder /app/.env /.env
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/feedbackboard"]

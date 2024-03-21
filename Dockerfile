FROM golang:1.21-alpine as builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o store ./cmd

FROM busybox as runner
COPY --from=builder /app/store /app
COPY config.yaml /etc/dbcore/config.yaml
ENTRYPOINT ["/app"]
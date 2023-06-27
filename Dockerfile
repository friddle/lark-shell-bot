FROM golang:1.18 as builder

ENV GOPROXY https://goproxy.cn

# mod
WORKDIR /app
COPY . /app/
RUN go mod download
RUN go build -o dist/feishu_shell_bot ./main.go

FROM alpine:latest
USER root

COPY  conf   /app/conf
COPY --from=builder /app/feishu_shell_bot   /app/feishu_shell_bot
WORKDIR /app/

ENTRYPOINT [ "/app/feishu_shell_bot" ]

FROM golang:1.18 as builder

ENV GOPROXY https://goproxy.cn

# mod
WORKDIR /app
COPY . /app/
RUN go mod download
RUN go build -o dist/feishu_shell_bot ./main.go

FROM golang:1.18
USER root
COPY --from=builder /app/dist/feishu_shell_bot /app/feishu_shell_bot
RUN chmod +x /app/feishu_shell_bot
ENV FEISHU_APP_ID=cli_aaaaaaaaaaaaa
ENV FEISHU_APP_SECRET=qmaaaaaaaaaaaaaaa
ENV FEISHU_ENCRYPT_KEY=mAAAAAAAAAAAAAAAAA
ENV FEISHU_VERIFICATION_TOKEN=HAAAAAAAAAAAAAAA
ENV FEISHU_BOT_PATH=/
ENV FEISHU_BOT_PORT=8080
WORKDIR /app/

ENTRYPOINT [ "/app/feishu_shell_bot" ]

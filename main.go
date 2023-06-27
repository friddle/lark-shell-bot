package main

import (
	"fmt"
	"github.com/go-zoox/logger"
	"lark-shell-bot/src"
)

func main() {
	feishuConfig := src.ReadFeishuConfig()
	println(fmt.Sprintf("info:%v", feishuConfig))
	bot, err := src.FeishuServer(feishuConfig)
	if err != nil {
		logger.Fatalf("bot error:%v", err)
	}
	if err := bot.Run(); err != nil {
		logger.Fatalf("bot error:%v", err)
	}
}

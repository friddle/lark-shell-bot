package main

import (
	"fmt"
	"lark-shell-bot/src"
)

func main() {
	feishuConfig := src.ReadFeishuConfig()
	println(fmt.Sprintf("info:%v", feishuConfig))
	src.FeishuServer(feishuConfig)
}

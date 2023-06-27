package src

import (
	"github.com/go-zoox/chatbot-feishu"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func ReadFeishuConfig() *chatbot.Config {
	//假如有文件。读取文件
	if Exists(".feishu.env") {
		err := godotenv.Load(".feishu.env")
		if err != nil {
			logs.Fatalf("read err %v", err)
		}
		logs.Println("load local feishu.env files")
	}

	port, _ := strconv.Atoi(os.Getenv("FEISHU_BOT_PORT"))
	//然后判断文件
	conf := chatbot.Config{
		AppID:             os.Getenv("FEISHU_APP_ID"),
		AppSecret:         os.Getenv("FEISHU_APP_SECRET"),
		EncryptKey:        os.Getenv("FEISHU_ENCRYPT_KEY"),
		VerificationToken: os.Getenv("FEISHU_VERIFICATION_TOKEN"),
		Port:              int64(port),
		Path:              os.Getenv("FEISHU_BOT_PATH"),
	}
	if conf.AppSecret == "" || conf.AppSecret == "" {
		logs.Fatalln("请配置APP_ID和APP_SECRET")
		os.Exit(2)
	}
	logs.Println("配置读取成功 %v", conf)
	return &conf
}

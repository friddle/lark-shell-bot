package src

import (
	"github.com/go-zoox/chatbot-feishu"
	chatgpt "github.com/go-zoox/chatgpt-client"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type SSHServerConfig struct {
	Hostname  string `yaml:"hostname"`
	IPAddress string `yaml:"ipaddress"`
	Password  string `yaml:"password"`
	Port      string `yaml:"port"`
	Username  string `yaml:"username"`
}

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
	if conf.AppID == "" || conf.AppSecret == "" {
		logs.Fatalln("请配置APP_ID和APP_SECRET")
		os.Exit(2)
	}
	logs.Printf("配置读取成功 %v", conf)
	return &conf
}

func ReadSshServer() map[string]*SSHServerConfig {
	configs := make(map[string]*SSHServerConfig, 0)
	if Exists(".machines.yaml") {
		err := ReadYamlFromFile(".machines.yaml", &configs)
		if err != nil {
			logs.Fatalf("read err %v", err)
		}
		logs.Println("load local machines.yaml files")
	}
	for name, config := range configs {
		config.Hostname = name
	}
	return configs
}

func ReadChatGptClient() *chatgpt.Config {
	if Exists(".chatgpt.env") {
		err := godotenv.Load(".chatgpt.env")
		if err != nil {
			logs.Fatalf("read err %v", err)
		}
		logs.Println("load local chatgpt.env files")
	}
	//然后判断文件
	conf := chatgpt.Config{
		APIKey:               os.Getenv("CHATGPT_API_KEY"),
		APIServer:            os.Getenv("CHATGPT_API_SERVER"),
		APIType:              os.Getenv("CHATGPT_API_TYPE"),
		AzureResource:        os.Getenv("CHATGPT_AZURE_RESOURCE"),
		AzureDeployment:      os.Getenv("CHATGPT_AZURE_DEPLOYMENT"),
		AzureAPIVersion:      os.Getenv("CHATGPT_AZURE_API_VERSION"),
		ConversationContext:  os.Getenv("CHATGPT_CONVERSATION_CONTEXT"),
		ConversationLanguage: os.Getenv("CHATGPT_CONVERSATION_LANGUAGE"),
		ChatGPTName:          os.Getenv("CHATGPT_NAME"),
		Proxy:                os.Getenv("CHATGPT_PROXY"),
	}
	if conf.APIKey == "" {
		return nil
	}
	return &conf
}

package src

import (
	"errors"
	"fmt"
	chatgpt "github.com/go-zoox/chatgpt-client"
	"github.com/go-zoox/uuid"
	"log"
	"os"
	regexp "regexp"
	"strings"
)

type SSHGptClient struct {
	client     chatgpt.Client
	config     chatgpt.Config
	promptText string
}

func NewChatGptClient() *SSHGptClient {
	config := ReadChatGptClient()
	if config == nil {
		return nil
	}
	gptClient, err := chatgpt.New(config)
	if err != nil {
		return nil
	}
	var promptBytes []byte
	if Exists("Prompt.txt") {
		promptBytes, _ = os.ReadFile("Prompt.txt")
	} else {
		promptBytes, _ = PromptFs.ReadFile("prompt/Prompt.txt")
	}
	log.Printf("prompt 信息为:" + string(promptBytes))

	client := &SSHGptClient{
		client:     gptClient,
		promptText: string(promptBytes) + "\r 描述为: %s",
	}
	client.RunInit()
	return client
}

// 返回是 执行的command和 效果
func (client *SSHGptClient) TranslateChatgptCmd(id string, commandText ...string) (string, error) {
	info := strings.Join(commandText, "")
	askCommand := fmt.Sprintf(client.promptText, info)

	conversation, err := client.client.GetOrCreateConversation(uuid.V4(), &chatgpt.ConversationConfig{})
	if err != nil {
		return "", err
	}
	outputCmdBytes, err := conversation.Ask([]byte(askCommand), &chatgpt.ConversationAskConfig{
		ID: id,
	})
	if err != nil {
		return "", err
	}

	outputCmd := string(outputCmdBytes)
	if err != nil {
		return "", err
	}
	if outputCmd == "" {
		return "", errors.New("输入文本无法理解相应问题")
	}

	rx := regexp.MustCompile(`command.*\[(.*)\]`)
	// 在字符串中查找匹配项
	match := rx.FindStringSubmatch(outputCmd)
	if len(match) == 2 {
		outputCmd = match[1]
	} else {
		return "", errors.New(fmt.Sprintf("chatgpt返回的命令不对:%s", outputCmd))
	}
	if outputCmd == "" {
		return "", errors.New(fmt.Sprintf("chatgpt无法理解你的输入%s", commandText))
	}
	log.Printf("chatgpt ask:%s \r response:%s\r", askCommand, outputCmd)
	return outputCmd, nil
}

func (client *SSHGptClient) RunInit() error {
	return nil
}

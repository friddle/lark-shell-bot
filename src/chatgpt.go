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
	client       chatgpt.Client
	config       chatgpt.Config
	promptText   string
	describeText string
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

	client := &SSHGptClient{
		client:     gptClient,
		promptText: string(promptBytes),
		describeText: "请把如下的描述翻译成固定的Linux命令格式。" +
			"并以 command:[Linux命令]的格式输出,假如无法识别为Linux命令，则返回为\"command:[]\"\n" +
			"描述为:%s",
	}
	client.RunInit()
	return client
}

// 返回是 执行的command和 效果
func (client *SSHGptClient) TranslateChatgptCmd(commandText ...string) (string, error) {
	info := strings.Join(commandText, "")
	askCommand := fmt.Sprintf(client.describeText, info)
	log.Printf("chatgpt message %s", askCommand)

	conversation, err := client.client.GetOrCreateConversation(uuid.V4(), &chatgpt.ConversationConfig{})
	if err != nil {
		return "", err
	}
	outputCmdBytes, err := conversation.Ask([]byte(askCommand), &chatgpt.ConversationAskConfig{})
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
		return "", errors.New(fmt.Sprintf("chatgpt返回的命令不对 %s", outputCmd))
	}
	return outputCmd, nil
}

func (client *SSHGptClient) RunInit() error {
	return nil
}

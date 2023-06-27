package src

import (
	"regexp"
	"testing"
)

func TestChatGptServer(t *testing.T) {
	gptClient := NewChatGptClient()
	if gptClient == nil {
		logs.Fatal("gpt client is fatal nil")
	}
	info, err := gptClient.TranslateChatgptCmd("查看目录")
	logs.Println(info, err)
}

func TestServerCmd(t *testing.T) {
	serverCmd := NewServerCmdClient()
	info, err := serverCmd.executeCmd("machine-1", "ls")
	logs.Println(info, err)
}

func TestRegex(t *testing.T) {
	outputCmd := "命令:[docker images | grep nginx]"
	rx := regexp.MustCompile(`命令:\[(.*)\]`)
	// 在字符串中查找匹配项
	match := rx.FindStringSubmatch(outputCmd)
	if len(match) == 2 {
		outputCmd = match[1]
	} else {
		logs.Fatal("error")
	}

}

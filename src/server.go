package src

import (
	"bufio"
	"fmt"
	"github.com/go-zoox/chatbot-feishu"
	"github.com/go-zoox/core-utils/regexp"
	"github.com/go-zoox/feishu"
	"github.com/go-zoox/feishu/contact/user"
	"github.com/go-zoox/feishu/event"
	feishuEvent "github.com/go-zoox/feishu/event"
	mc "github.com/go-zoox/feishu/message/content"
	"github.com/go-zoox/logger"
	"io"
	"os/exec"
	"strings"
	"time"
)

// 两种。一种是快速的..一种是长时间执行的

func getUser(request *feishuEvent.EventRequest) (*user.RetrieveResponse, error) {
	sender := request.Sender()
	return &user.RetrieveResponse{
		User: user.UserEntity{
			Name:    sender.SenderID.UserID,
			OpenID:  sender.SenderID.OpenID,
			UnionID: sender.SenderID.UnionID,
			UserID:  sender.SenderID.UserID,
		},
	}, nil
}

func ReplyText(reply func(context string, msgType ...string) error, text string) error {
	msgType, content, err := mc.
		NewContent().
		Post(&mc.ContentTypePost{
			ZhCN: &mc.ContentTypePostBody{
				Content: [][]mc.ContentTypePostBodyItem{
					{
						{
							Tag:      "text",
							UnEscape: false,
							Text:     text,
						},
					},
				},
			},
		}).
		Build()
	if err != nil {
		return fmt.Errorf("failed to build content: %v", err)
	}
	if err := reply(string(content), msgType); err != nil {
		logs.Printf("content :%s", content)
		logger.Info(fmt.Sprintf("failed to reply: %v", err))
	}

	return nil
}

func getCommand(client feishu.Client, text string, request *feishuEvent.EventRequest) string {
	var command string
	// group chat
	if request.IsGroupChat() {
		botInfo, _ := client.Bot().GetBotInfo()
		if ok := regexp.Match("^@_user_1", text); ok {
			for _, mention := range request.Event.Message.Mentions {
				if mention.Key == "@_user_1" && mention.ID.OpenID == botInfo.OpenID {
					command = strings.Replace(text, "@_user_1", "", 1)
					logger.Info("chat command %s", command)
					break
				}
			}
		}
	} else if request.IsP2pChat() {
		command = text
	}
	command = strings.TrimSpace(command)
	logger.Info("chat command %s", command)
	return command
}

func FeishuServer(feishuConf *chatbot.Config) (chatbot.ChatBot, error) {
	bot, err := chatbot.New(feishuConf)
	client := feishu.New(&feishu.Config{"https://open.feishu.cn", feishuConf.AppID, feishuConf.AppSecret})
	chatGptClient := NewChatGptClient()
	sshServerClient := NewServerCmdClient()
	if err != nil {
		logger.Errorf("failed to create bot: %v", err)
		return nil, err
	}

	bot.OnCommand("ping", &chatbot.Command{
		Handler: func(args []string, request *feishuEvent.EventRequest, reply func(content string, msgType ...string) error) error {
			if err := ReplyText(reply, "pong"); err != nil {
				return fmt.Errorf("failed to reply: %v", err)
			}
			return nil
		},
	})

	bot.OnCommand("help", &chatbot.Command{
		Handler: func(args []string, request *event.EventRequest, reply chatbot.MessageReply) error {
			if err := ReplyText(reply, "run command bro"); err != nil {
				return fmt.Errorf("failed to reply: %v", err)
			}
			return nil
		},
	})

	bot.OnCommand("machines", &chatbot.Command{
		Handler: func(args []string, request *event.EventRequest, reply chatbot.MessageReply) error {
			var machineTexts []string
			for _, machine := range sshServerClient.machines {
				machineTexts = append(machineTexts, machine.Hostname+fmt.Sprintf("(%s)", machine.IPAddress))
			}
			ReplyText(reply, strings.Join(machineTexts, "\r\n"))
			return nil
		},
	})

	bot.OnMessage(func(text string, request *event.EventRequest, reply chatbot.MessageReply) error {
		command := getCommand(client, text, request)
		commands := strings.Split(command, " ")
		if command == "" {
			logger.Infof("ignore empty command message")
			return nil
		}
		if strings.HasPrefix(command, "/chatgpt") {
			if chatGptClient == nil {
				ReplyText(reply, "ChatGpt 没有设置或者相应配置不正确")
				return nil
			}
			command, err := chatGptClient.TranslateChatgptCmd(commands[1:]...)
			if err == nil {
				ReplyText(reply, fmt.Sprintf("执行command:%s", command))
				RunCommand(reply, command)
			} else {
				ReplyText(reply, fmt.Sprintf("执行命令失败 %v", err))
			}
			return nil
		}
		if strings.HasPrefix(command, "/ssh") {
			output, err := sshServerClient.executeCmd(commands[1], commands[2:]...)
			if err != nil {
				ReplyText(reply, fmt.Sprintf("err %v", err))
			}
			ReplyText(reply, output)
			return nil
		}

		if strings.HasPrefix(command, "/") {
			logger.Infof("ignore empty command message")
			return nil
		}
		go func() {
			RunCommand(reply, command)
		}()
		return nil
	})

	return bot, nil
}

func RunCommand(reply chatbot.MessageReply, command string) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("%s", command))
	logs.Println(cmd)
	stderr, err := cmd.StderrPipe()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		ReplyText(reply, fmt.Sprintf("run command error %v", err))
	}
	// Start command
	if err := cmd.Start(); err != nil {
		logger.Infof("error starting command:%v", err)
		ReplyText(reply, fmt.Sprintf("error run command bro:%v", err))
		return
	}

	stdall := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(stdall)
	i := 0
	for {
		texts := make([]string, 0)
		for scanner.Scan() {
			texts = append(texts, scanner.Text())
		}
		if len(texts) == 0 {
			break
		}
		ReplyText(reply, strings.Join(texts, "\r\n"))
		time.Sleep(1 * time.Second)
		i = i + 1
		if i > 20 {
			ReplyText(reply, "命令运行太久。直接退出")
			break
		}
	}
	if err != nil {
		ReplyText(reply, fmt.Sprintf("run command error %v", err))
	}
	defer stdout.Close()
	defer stderr.Close()
}

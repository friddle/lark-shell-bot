package src

import (
	"bufio"
	"fmt"
	"github.com/go-zoox/chatbot-feishu"
	"github.com/go-zoox/feishu/contact/user"
	"github.com/go-zoox/feishu/event"
	feishuEvent "github.com/go-zoox/feishu/event"
	mc "github.com/go-zoox/feishu/message/content"
	"github.com/go-zoox/logger"
	"io"
	"os/exec"
	"strings"
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
							UnEscape: true,
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
		return fmt.Errorf("failed to reply: %v", err)
	}

	return nil
}

func getCommand(text string, request *feishuEvent.EventRequest) string {
	var command string
	// group chat
	if request.IsGroupChat() {
		command = text
	} else if request.IsP2pChat() {
		command = text
	}
	logger.Info("command %s", command)
	return command
}

func FeishuServer(feishuConf *chatbot.Config) {
	bot, err := chatbot.New(feishuConf)
	if err != nil {
		logger.Errorf("failed to create bot: %v", err)
		return
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

	bot.OnCommand("chatgpt", &chatbot.Command{
		Handler: func(args []string, request *event.EventRequest, reply chatbot.MessageReply) error {
			if err := ReplyText(reply, "没有实现Gpt哦bro"); err != nil {
				return fmt.Errorf("failed to reply: %v", err)
			}
			return nil
		},
	})

	bot.OnCommand("list", &chatbot.Command{
		Handler: func(args []string, request *event.EventRequest, reply chatbot.MessageReply) error {
			RunCommand(reply, "ls")
			return nil
		},
	})

	bot.OnMessage(func(text string, request *event.EventRequest, reply chatbot.MessageReply) error {
		command := getCommand(text, request)
		if command == "" {
			logger.Infof("ignore empty command message")
			return nil
		}
		if strings.HasPrefix(command, "/") {
			logger.Infof("ignore empty command message")
			return nil
		}
		cmdArr := strings.Split(command, " ")
		cmd := cmdArr[0]
		args := cmdArr[1:]

		RunCommand(reply, cmd, args...)
		return nil
	})
}

func RunCommand(reply chatbot.MessageReply, command string, args ...string) {
	cmd := exec.Command(command, args...)
	stderr, err := cmd.StderrPipe()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err := ReplyText(reply, fmt.Sprintf("run command error %v", err))
		if err != nil {
			logger.Info("error reply message")
		}
	}
	stdall := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(stdall)
	for scanner.Scan() {
		err := ReplyText(reply, scanner.Text())
		if err != nil {
			logger.Info("error reply message")
		}
	}
	err = cmd.Wait()
	if err != nil {
		err := ReplyText(reply, fmt.Sprintf("run command error %v", err))
		if err != nil {
			logger.Info("error reply message")
		}
	}
	defer stdout.Close()
	defer stderr.Close()
}

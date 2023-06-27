package src

import (
	"errors"
	"fmt"
	"strings"
)

type ServerCmdClient struct {
	machines map[string]*SSHServerConfig
}

func NewServerCmdClient() *ServerCmdClient {
	configs := ReadSshServer()
	client := ServerCmdClient{
		machines: configs,
	}
	return &client
}

func (client *ServerCmdClient) executeCmd(machine string, cmd ...string) (string, error) {
	if machine, ok := client.machines[machine]; ok {
		sshClient, err := NewSSHClient(machine.IPAddress, machine.Port, machine.Username, machine.Password)
		if err != nil {
			return "", err
		}
		smartOutput, err := sshClient.Sshclient.Script(strings.Join(cmd, " ")).SmartOutput()
		return string(smartOutput), err
	} else {
		return "", errors.New(fmt.Sprintf("没有找到服务器%s", machine))
	}
}

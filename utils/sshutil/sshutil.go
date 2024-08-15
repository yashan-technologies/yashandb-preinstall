package sshutil

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type SSH struct {
	IP       string
	User     string
	Password string
	Port     int
}

func NewSSH(ip string, port int, user, password string) *SSH {
	return &SSH{
		IP:       ip,
		User:     user,
		Password: password,
		Port:     port,
	}
}

func (s SSH) CheckSSHConnection() error {
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 注意：这会忽略SSH密钥验证，仅用于测试
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port), config)
	if err != nil {
		return err
	}
	defer client.Close()
	return nil
}

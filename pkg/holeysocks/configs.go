package holeysocks

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

// MainConfig contains SSH and Socks configuration variables
type MainConfig struct {
	SSH   sshConfig   `json:"ssh"`
	Socks socksConfig `json:"socks"`
}

type sshConfig struct {
	Username string `json:"username"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	PrivKey  []ssh.AuthMethod
}

func (s *sshConfig) setKey(keyBytes []byte) error {
	privateKey, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return err
	}
	auth := ssh.PublicKeys(privateKey)
	s.PrivKey = []ssh.AuthMethod{auth}
	return err
}

func (s *sshConfig) connectionString() string {
	return fmt.Sprintf("%s:%v", s.Host, s.Port)
}

type socksConfig struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

package holeysocks

import (
	"fmt"

	"github.com/armon/go-socks5"
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

func (s *sshConfig) SetKey(keyBytes []byte) error {
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
	Remote string `json:"remote"`
}

// ForwardService implements reverse port forwarding, similar to the -R flag
// in openssh-client. Configuration is done in the configs/ssh.json file.
// NOTE The generated keys and config.json data are embedded in the binary so
// take the appropriate precautions when setting up the ssh server user.
func ForwardService(config MainConfig) error {
	sshClientConf := &ssh.ClientConfig{
		User:            config.SSH.Username,
		Auth:            config.SSH.PrivKey,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH server
	serverConn, err := ssh.Dial("tcp", config.SSH.connectionString(), sshClientConf)
	if err != nil {
		return fmt.Errorf("dial INTO remote server error: %s", err)
	}

	// Create a listening port on the remote SSH server
	remoteListener, err := serverConn.Listen("tcp", config.Socks.Remote)
	if err != nil {
		return fmt.Errorf("unable to bind on remote's port: %s", err)
	}
	defer remoteListener.Close()

	// Create a net.Conn to the remote port and wait for a connection
	socksClient, err := remoteListener.Accept()
	if err != nil {
		return fmt.Errorf("could not handle incoming socks client request: %s", err)
	}

	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		return fmt.Errorf("failed creating new Socks server: %s", err)
	}
	return server.ServeConn(socksClient)
}

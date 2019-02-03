package holeysocks

import (
	"fmt"
	"io"
	"log"
	"net"

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
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

// Handle local client connections and tunnel data to the remote server
// will use io.Copy - http://golang.org/pkg/io/#Copy
// https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Fatalln(fmt.Sprintf("error while copy remote->local: %s", err))
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Fatalln(fmt.Sprintf("error while copy local->remote: %s", err))
		}
		chDone <- true
	}()

	<-chDone
}

// ForwardService implements reverse port forwarding similar to the -R flag
// in openssh-client. Configuration is done in the /configs/config.json file.
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
		return fmt.Errorf("Dial INTO remote server error: %s", err)
	}

	// Publish the designated local port to the same port on the remote SSH server
	remoteListener, err := serverConn.Listen("tcp", config.Socks.Remote)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("INFO: %s", err))
	}
	defer remoteListener.Close()

	// Handle incoming requests from the remote tunnel
	for {
		// Grab a handle to the pre-configured local port that will be sent to the remote
		// SSH server
		local, err := net.Dial("tcp", config.Socks.Local)
		if err != nil {
			return fmt.Errorf("Unable to start local listen: %s", err)
		}

		// Grab a handle on the remote port
		remote, err := remoteListener.Accept()
		if err != nil {
			return fmt.Errorf("Unable to accept remote traffic locally: %s", err)
		}

		// Swap IO from the local and remote hanles
		handleClient(remote, local)
	}
}

// DarnSocks creates a new SOCKS5 server at the provided ports and
// remote-forwards the port to another machine over SSH
func DarnSocks(config MainConfig) error {
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		return err
	}

	go func() {
		// Create a SOCKS5 server
		err = server.ListenAndServe("tcp", config.Socks.Local)
		if err != nil {
			fmt.Printf("ERR: %s", err.Error())
		}
	}()

	go func() {
		// Publish SOCKS to remote server
		// TODO make this goroutine wait for SOCKS
		err = ForwardService(config)
		if err != nil {
			fmt.Printf("ERR: %s", err.Error())
		}
	}()
	return err
}

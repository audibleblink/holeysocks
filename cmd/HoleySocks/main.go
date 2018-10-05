package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/armon/go-socks5"
	"github.com/gobuffalo/packr"
	"golang.org/x/crypto/ssh"
)

type mainConfig struct {
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

var config = mainConfig{}

func init() {
	// unpack the configs and ssh keys from the binary
	// the were packed at compile-time
	box := packr.NewBox("../../configs")
	configBytes := box.Bytes("config.json")
	if err := json.Unmarshal(configBytes, &config); err != nil {
		panic(err)
	}

	privKeyBytes := box.Bytes("id_ed25519")
	config.SSH.setKey(privKeyBytes)
}

// forwardService implements reverse port forwarding similar to the -R flag
// in openssh-client. Configuration is done in the /configs/ssh.json file.
// NOTE The generated keys and ssh.json data are embedded in the binary so
// take the appropriate precautions when setting up the ssh server user.
func forwardService() {

	sshClientConf := &ssh.ClientConfig{
		User:            config.SSH.Username,
		Auth:            config.SSH.PrivKey,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH server
	serverConn, err := ssh.Dial("tcp", config.SSH.connectionString(), sshClientConf)
	if err != nil {
		log.Fatalln(fmt.Sprintf("SSH Conn Failed: %s", err))
	}

	// Publish the designated local port to the same port on the remote SSH server
	remoteListener, err := serverConn.Listen("tcp", config.Socks.Remote)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Reverse port-forward failed : %s", err))
	}
	defer remoteListener.Close()

	// Handle incoming request from the remote tunnel
	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded
		local, err := net.Dial("tcp", config.Socks.Local)
		if err != nil {
			log.Fatalln(fmt.Sprintf("Unable to start local listen: %s", err))
		}

		remote, err := remoteListener.Accept()
		if err != nil {
			log.Fatalln(fmt.Sprintf("Unable to accept remote traffic locally: %s", err))
		}

		handleClient(remote, local)
	}

}

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote server
// Will use io.Copy - http://golang.org/pkg/io/#Copy
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

func darnSocks() {
	// Create a SOCKS5 server
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost
	server.ListenAndServe("tcp", config.Socks.Local)
}

func main() {
	go darnSocks()
	fmt.Printf("Serving on remote: %s", config.Socks.Remote)
	forwardService()
}

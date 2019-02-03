package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/audibleblink/HoleySocks/pkg/holeysocks"
	"github.com/gobuffalo/packr"
)

var (
	static    string
	sshUser   string
	sshHost   string
	sshPort   int
	rPort     int
	socksPort int
	pKey      string
)

func init() {
	flag.StringVar(&sshUser, "sshuser", "", "[REQ] SSH user ong the host")
	flag.StringVar(&sshHost, "sshhost", "", "[REQ] SSH host with which to connect")
	flag.StringVar(&pKey, "pkey", "", "[REQ] File path for private key")
	flag.IntVar(&sshPort, "sshport", 22, "SSH host destination port")
	flag.IntVar(&rPort, "rport", 1080, "SSH host port on which to bind the local SOCKS server")
	flag.IntVar(&socksPort, "socksport", 1080, "Bind port of the SOCKS server")
	flag.Parse()
}

func main() {
	config := holeysocks.MainConfig{}

	switch static {
	case "1":
		box := packr.NewBox("../../configs")
		configBytes, err := box.Find("ssh.json")
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(configBytes, &config)
		if err != nil {
			panic(err)
		}

		privKeyBytes, err := box.Find("id_ed25519")
		if err != nil {
			panic(err)
		}
		config.SSH.SetKey(privKeyBytes)
	default:
		if sshUser == "" || sshHost == "" || pKey == "" {
			panic("Missing required flag. Use -help")
		}
		config.SSH.Username = sshUser
		config.SSH.Host = sshHost
		config.SSH.Port = sshPort
		config.Socks.Local = fmt.Sprintf("127.0.0.1:%d", socksPort)
		config.Socks.Remote = fmt.Sprintf("127.0.0.1:%d", rPort)

		pKeyBytes, err := ioutil.ReadFile(pKey)
		if err != nil {
			panic(err)
		}

		err = config.SSH.SetKey(pKeyBytes)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Serving on remote: %s\n", config.Socks.Remote)
	err := holeysocks.DarnSocks(config)
	if err != nil {
		panic(err)
	}

	// DarnSocks is concurrent, so we must keep main from exiting
	select {}
}

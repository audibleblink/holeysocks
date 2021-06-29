package main

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/audibleblink/holeysocks/pkg/holeysocks"
)

//go:embed configs/ssh.json
var configBytes []byte

//go:embed configs/id_ed25519
var pKeyBytes []byte

func main() {
	config := holeysocks.MainConfig{}

	err := json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	err = config.SSH.SetKey(pKeyBytes)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Serving on remote: %s\n", config.Socks.Remote)
	err = holeysocks.ForwardService(config)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"github.com/audibleblink/HoleySocks/pkg/holeysocks"
)

func main() {
	fmt.Printf("Serving on remote: %s\n", holeysocks.Config.Socks.Remote)
	err := holeysocks.DarnSocks()
	if err != nil {
		fmt.Printf("ERR: %s", err.Error())
	}
}

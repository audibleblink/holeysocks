# HoleySocks

A simple cross-platform reverse socks proxy.


## Getting Started

### As a module

```go
import github.com/audibleblink/HoleySocks/pkg/holeysocks

func main() {
	config := holeysocks.MainConfig{}
        configBytes, _ := ioutil.OpenFile("ssh.json")
        json.Unmarshal(configBytes, &config)
        holeysocks.DarnSocks(config)
}
```


### As a standalone binary

It's possible to embed all the required parameters to start and forward
the socks server with SSH so that cli flags are not needed.
Do this by creating `config/ssh.json` and using the `-X main.static=1` ldflag.

```bash
# needed for embedding configs in the binary
go get -u github.com/gobuffalo/packr/...

go get github.com/audibleblink/HoleySocks/...
cd $GOPATH/src/github.com/audibleblink/HoleySocks

... edit configs/ssh.json ...

make depends
make
```

To compile a generic binary without embedded configs, remote the `-X` ldflag from the `Makefile` or 
just `go build` as necessary. You should get a binary that's configurable with these flags:

```
Usage of binaries/linux/HoleySocks64:
  -sshuser string
        [REQ] SSH user ong the host
  -sshhost string
        [REQ] SSH host with which to connect
  -pkey string
        [REQ] File path for private key
  -rport int
        SSH host port on which to bind the local SOCKS server (default 1080)
  -socksport int
        Bind port of the SOCKS server (default 1080)
  -sshport int
        SSH host destination port (default 22)
```

Read the Makefile for more options

**CAUTION**
The generated private keys are embedded into the binary to allow for the reverse
port forwarding without interaction. Follow the instructions below.

Before running the generated binaries, you'll need a user on your attacking machine
for receiving the reverse ssh connection that forwards the socks proxy from the victim.

Once that user has been created, (with a homedir and /bin/false shell), append the generated
pubkey in your authorized_keys file on the attacking machine.

Do so with the following prefixes:

```
# if you're forwarding port 1080
FROM=<victim_ip_or_host> NO-X11-FORWARDING,PERMITOPEN="0.0.0.0:1080" ssh-ed25519 AAAAC3......
```

The Makefile should generate the needed commands and entry for you when you run `make depends`

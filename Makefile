NAME=HoleySocks

OUT_LINUX=binaries/linux/${NAME}
OUT_MACOS=binaries/macos/${NAME}
OUT_WINDOWS=binaries/windows/${NAME}

BUILD=packr2 build
SRC=cmd/HoleySocks/*

STRIP=-s
LINUX_LDFLAGS=--ldflags "${STRIP} -w"
WIN_LDFLAGS=--ldflags "${STRIP} -w -H windowsgui"

all: linux64 windows64 macos64 linux32 macos32 windows32 

depends:
	ssh-keygen -t ed25519 -f configs/id_ed25519 -N '' -C ${NAME}
	@echo
	@echo "================================================="
	@echo "Create a user with a /bin/false shell on the target ssh server."
	@echo "useradd -s /bin/false -m -d /home/sshuser -N sshuser"
	@echo
	@echo "Append the following line to that user's authorized_keys file:"
	@echo "NO-X11-FORWARDING,PERMITOPEN=\"0.0.0.0:1080\" `cat ./configs/id_ed25519.pub`"
	@echo
	@echo "If you know your target's public IP, you can also prepend the above with:"
	@echo "FROM=<ip or hostname>"
	@echo "================================================="

linux64:
	GOOS=linux GOARCH=amd64 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX}64 ${SRC}
macos64:
	GOOS=darwin GOARCH=amd64 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_MACOS}64 ${SRC}
windows64:
	GOOS=windows GOARCH=amd64 ${BUILD} ${WIN_LDFLAGS} -o ${OUT_WINDOWS}64.exe ${SRC}
linux32:
	GOOS=linux GOARCH=386 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX}32 ${SRC}
macos32:
	GOOS=darwin GOARCH=386 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_MACOS}32 ${SRC}
windows32:
	GOOS=windows GOARCH=386 ${BUILD} ${WIN_LDFLAGS} -o ${OUT_WINDOWS}32.exe ${SRC}

.PHONY: linux64 windows64 macos64 linux32 macos32 windows32

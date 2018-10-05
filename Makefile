NAME=HoleySocks

OUT_LINUX=binaries/linux/${NAME}
OUT_MACOS=binaries/macos/${NAME}
OUT_WINDOWS=binaries/windows/${NAME}

BUILD=packr build
SRC=cmd/HoleySocks/main.go

STRIP=-s
LINUX_LDFLAGS=--ldflags "${STRIP} -w"
WIN_LDFLAGS=--ldflags "${STRIP} -w -H windowsgui"

MINGW=x86_64-w64-mingw32-gcc-7.3-posix

all: linux64 windows64 macos64 linux32 macos32 windows32 

depends:
	ssh-keygen -t ed25519 -f configs/id_ed25519 -N ''
	echo "Create a user with a /bin/false shell on the target ssh server that will be used "
	echo "for socks forwarding. Also append the genereate pubkey to authorized_keys"

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

.PHONY: linux64 windows64 macos64 linux32 macos32 windows32 clean listen

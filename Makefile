VERSION := v1.0
GITCOMMIT := $(shell git describe --always)
GITBRANCH := $(shell git rev-parse --abbrev-ref HEAD)
TIMESTAMP := $(shell date --iso=seconds)

.PHONY: workload-policy-manager installer docker all clean

workload-policy-manager:
	env GOOS=linux go build -ldflags "-X main.Version=$(VERSION)-$(GITCOMMIT) -X main.Branch=$(GITBRANCH) -X main.Time=$(TIMESTAMP)" -o out/workload-policy-manager

installer: workload-policy-manager
	mkdir -p out/wpm
	chmod +x dist/linux/build-secure-docker-daemon.sh
	dist/linux/build-secure-docker-daemon.sh
	cp -rf secure-docker-daemon/out out/wpm/docker-daemon
	rm -rf secure-docker-daemon
	cp -f dist/linux/daemon.json out/wpm/
	cp dist/linux/install.sh out/wpm/install.sh && chmod +x out/wpm/install.sh
	cp dist/linux/uninstall-secure-docker-daemon.sh out/wpm/uninstall-secure-docker-daemon.sh && chmod +x out/wpm/uninstall-secure-docker-daemon.sh
	cp out/workload-policy-manager out/wpm/workload-policy-manager
	chmod +x out/wpm/workload-policy-manager
	makeself out/wpm out/wpm-$(VERSION).bin "Workload Policy Manager $(VERSION)" ./install.sh

all: installer

clean:
	rm -rf out/

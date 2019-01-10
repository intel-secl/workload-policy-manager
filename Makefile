GITTAG := $(shell git describe --tags --abbrev=0 2> /dev/null)
GITCOMMIT := $(shell git describe --always)
GITCOMMITDATE := $(shell git log -1 --date=short --pretty=format:%cd)
VERSION := $(or ${GITTAG}, v0.0.0)

.PHONY: workload-policy-manager installer docker all clean

workload-policy-manager:
	env GOOS=linux go build -ldflags "-X intel/isecl/workload-policy-manager/version.Version=$(VERSION)-$(GITCOMMIT)" -o out/workload-policy-manager

installer: workload-policy-manager
	mkdir -p out/wpm
	cp dist/linux/install.sh out/wpm/install.sh && chmod +x out/wpm/install.sh
	cp out/workload-policy-manager out/wpm/workload-policy-manager
	makeself out/wpm out/wpm-$(VERSION).bin "Workload Policy Manager $(VERSION)" ./install.sh

all: installer

clean:
	rm -rf out/
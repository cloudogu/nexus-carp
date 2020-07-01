ARTIFACT_ID=nexus-carp
VERSION=1.0.0

TARGETDIR=target
PKG=${ARTIFACT_ID}-${VERSION}.tar.gz
BINARY=${TARGETDIR}/${ARTIFACT_ID}

default: build

include build/make/variables.mk
include build/make/self-update.mk
include build/make/clean.mk
include build/make/dependencies-gomod.mk

generate:
	go generate

setup: generate dependencies

$(BINARY): setup
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-X main.Version=${VERSION} -extldflags "-static"' -o $(BINARY) .

build: $(BINARY)

package: build
	cd ${TARGETDIR}; tar cvfz ${PKG} ${ARTIFACT_ID}

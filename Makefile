ARTIFACT_ID=nexus-carp
VERSION=1.5.0

MAKEFILES_VERSION=9.2.1

TARGETDIR=target
PKG=${ARTIFACT_ID}-${VERSION}.tar.gz
BINARY=${TARGETDIR}/${ARTIFACT_ID}

default: build

include build/make/variables.mk
include build/make/self-update.mk
include build/make/clean.mk
include build/make/dependencies-gomod.mk
include build/make/digital-signature.mk

generate:
	go generate

setup: generate dependencies

$(BINARY): setup
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-X main.Version=${VERSION} -extldflags "-static"' -o $(BINARY) .

build: $(BINARY) signature

package: build signature
	cd ${TARGETDIR}; tar cvfz ${PKG} ${ARTIFACT_ID}

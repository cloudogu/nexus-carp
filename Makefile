APP=nexus-carp
VERSION=0.3.2

TARGETDIR=target
PKG=${APP}-${VERSION}.tar.gz
BINARY=${TARGETDIR}/${APP}

default: build

dependencies:
	dep ensure

generate:
	go generate

setup: dependencies generate

$(BINARY): setup
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-X main.Version=${VERSION} -extldflags "-static"' -o $(BINARY) .

build: $(BINARY)

package: build
	cd ${TARGETDIR}; tar cvfz ${PKG} ${APP}

clean:
	rm -rf $(TARGETDIR)

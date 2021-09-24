PREFIX        = /usr
GO_BUILDFLAGS = -mod vendor
GO_LDFLAGS    = 

all: FORCE
	go build $(GO_BUILDFLAGS) -ldflags '-s -w $(GO_LDFLAGS)' -o build/art .

install: FORCE all
	install -D -m 0755 build/art "$(DESTDIR)$(PREFIX)/bin/art"

vendor: FORCE
	go mod tidy
	go mod vendor
	go mod verify

.PHONY: FORCE

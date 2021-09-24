PREFIX = /usr

GO            = GOBIN=$(CURDIR)/build go
GO_BUILDFLAGS = -mod vendor
GO_LDFLAGS    = -s -w

all: FORCE
	$(GO) install $(GO_BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' .

install: FORCE all
	install -D -m 0755 build/art "$(DESTDIR)$(PREFIX)/bin/art"

vendor: FORCE
	go mod tidy
	go mod vendor
	go mod verify

.PHONY: FORCE

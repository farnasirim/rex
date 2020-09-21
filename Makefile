GO ?= go
REX := github.com/farnasirim/rex

.PHONY: rex rexd
all: rex rexd

rex rexd:
	$(GO) build $(REX)/cmd/$@

clean:
	rm -f rex rexd

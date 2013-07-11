GOROOT ?=	/usr/local/go
GO ?=		$(GOROOT)/bin/go

GOPATH :=	$(shell pwd)

all:
	GOPATH="$(GOPATH)" $(GO) install -v dumper

format:
	find . -name \*.go -exec go fmt {} \;

clean:
	find . -name \*~ -exec rm -f {} \;
	rm -rf bin pkg


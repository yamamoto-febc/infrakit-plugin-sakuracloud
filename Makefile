TEST?=$$(go list ./... | grep -v vendor)
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: test fmt

clean:
	rm -Rf $(CURDIR)/bin/*

build: clean fmt
	govendor build -ldflags "-s -w" -o $(CURDIR)/bin/sakuracloud $(CURDIR)/*.go

build-x: clean vet
	sh -c "'$(CURDIR)/scripts/build.sh'"

test: fmt
	govendor test $(TEST) $(TESTARGS) -v -timeout=30m -parallel=4 ;

fmt:
	gofmt -s -l -w $(GOFMT_FILES)

.PHONY: default test vet fmt lint

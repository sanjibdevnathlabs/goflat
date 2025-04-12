.PHONY: lint-pkgs lint test

lint-pkgs:
	go get -u honnef.co/go/tools/cmd/staticcheck
	go get -u github.com/client9/misspell/cmd/misspell

lint:
	$(exit $(go fmt ./... | wc -l))
	go vet ./...
	find . -type f -name "*.go" | xargs misspell -error -locale US
	staticcheck $(go list ./...)

test:
	go test -race ./...

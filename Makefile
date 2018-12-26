all: test

testdeps:
	go get -d -t ./...

checkfmt: testdeps
	[ -z "$$(gofmt -s -d . | tee /dev/stderr)" ]

lint: testdeps
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run -D errcheck -E golint -E staticcheck -E misspell -E gofmt

coverage: lint
	go test -coverprofile=coverage.txt -covermode=atomic ./...

test: lint
	go test ./...

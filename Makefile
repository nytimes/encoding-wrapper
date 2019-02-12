all: test

testdeps:
	go get -d -t ./...

checkfmt: testdeps
	[ -z "$$(gofmt -s -d . | tee /dev/stderr)" ]

lint: testdeps
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --enable-all -D lll -D errcheck -D dupl -D gochecknoglobals --deadline 5m ./...

coverage: lint
	go test -coverprofile=coverage.txt -covermode=atomic ./...

test: lint
	go test ./...

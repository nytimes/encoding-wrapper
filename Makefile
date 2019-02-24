all: test

testdeps:
	go mod download

checkfmt: testdeps
	[ -z "$$(gofmt -s -d . | tee /dev/stderr)" ]

lintdeps: testdeps
	GO111MODULE=off go get github.com/golangci/golangci-lint/cmd/golangci-lint

run-lint:
	golangci-lint run --enable-all -D lll -D errcheck -D dupl -D gochecknoglobals --deadline 5m ./...

lint: lintdeps run-lint

coverage: run-lint
	go test -coverprofile=coverage.txt -covermode=atomic ./...

test: run-lint
	go test ./...

all: test

testdeps:
	go get -d -t ./...

checkfmt: testdeps
	[ -z "$$(gofmt -s -d . | tee /dev/stderr)" ]

lint: testdeps
	go get github.com/alecthomas/gometalinter honnef.co/go/unused/cmd/unused
	gometalinter --install --vendored-linters
	go install ./...
	go list -f '{{.TestImports}}' ./... | sed -e 's/\[\(.*\)\]/\1/' | tr ' ' '\n' | grep '^.*\..*/.*$$' | xargs go install
	gometalinter -j 4 --enable=gofmt --enable=unused --disable=dupl --disable=errcheck --disable=gas --deadline=10m --tests ./...

test: lint
	go test ./...

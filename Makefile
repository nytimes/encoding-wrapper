all: test

testdeps:
	go get -d -t ./...

checkfmt: testdeps
	[ -z "$$(gofmt -s -d . | tee /dev/stderr)" ]

lint: testdeps
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run -D errcheck -E golint -E staticcheck -E misspell -E gofmt

coverage: lint
	@rm -f coverage.txt; for p in $$(go list ./...); do \
		go test -coverprofile=profile.out -covermode=atomic $$p || export status=2; \
		if [ -f profile.out ]; then cat profile.out >> coverage.txt; rm profile.out; fi \
		done; \
		exit $${status:-0}

test: lint
	go test ./...

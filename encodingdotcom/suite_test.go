package encodingdotcom

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/check.v1"
)

type S struct{}

func Test(t *testing.T) {
	check.TestingT(t)
}

var _ = check.Suite(&S{})

func (s *S) startServer(content string, status int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(content))
	}))
	return server
}

package elementalconductor

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type fakeServerRequest struct {
	req  *http.Request
	body []byte
}

func startServer(status int, content string) (*httptest.Server, chan fakeServerRequest) {
	requests := make(chan fakeServerRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := ioutil.ReadAll(r.Body)
		fakeRequest := fakeServerRequest{req: r, body: data}
		requests <- fakeRequest
		w.WriteHeader(status)
		w.Write([]byte(content))
	}))
	return server, requests
}

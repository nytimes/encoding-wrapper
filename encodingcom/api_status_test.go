package encodingcom

import (
	"net/http"
	"net/http/httptest"

	"gopkg.in/check.v1"
)

func (s *S) TestAPIStatus(c *check.C) {
	var req *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
		w.Write([]byte(`{"status":"Encoding Queue Processing Delays","status_code":"queue_slow","incident":"Our encoding queue is processing slower than normal.  Check back for updates."}`))
	}))
	defer server.Close()
	resp, err := APIStatus(server.URL)
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, APIStatusResponse{
		Status:     "Encoding Queue Processing Delays",
		StatusCode: "queue_slow",
		Incident:   "Our encoding queue is processing slower than normal.  Check back for updates.",
	})
}

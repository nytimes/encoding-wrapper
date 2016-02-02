package elementalcloud

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"gopkg.in/check.v1"
)

func (s *S) mockGenericResponseObject(message string, errors []string) interface{} {
	return map[string]interface{}{
		"response": map[string]interface{}{
			"message": message,
			"errors":  map[string][]string{"error": errors},
		},
	}
}

func (s *S) TestNewClient(c *check.C) {
	expected := Client{
		Host:    "https://mycluster.cloud.elementaltechnologies.com",
		UserID:  "myuser",
		UserKey: "secret-key",
	}
	got, err := NewClient("https://mycluster.cloud.elementaltechnologies.com", "myuser", "secret-key")
	c.Assert(err, check.IsNil)
	c.Assert(*got, check.DeepEquals, expected)
}

func (s *S) TestDoRequiredParameters(c *check.C) {
	var req *http.Request
	var data []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
		data, _ = ioutil.ReadAll(r.Body)
		w.Write([]byte(`<response>test</response>`))
	}))
	defer server.Close()
	client := Client{Host: server.URL, UserID: "myuser", UserKey: "123"}

	var respObj interface{}
	myJob := Job{FileInput: FileInput{URI: "some-file"}, Profile: "6"}
	err := client.do("POST", "/jobs", myJob, &respObj)

	c.Assert(err, check.IsNil)
	c.Assert(req, check.NotNil)
	c.Assert(req.Method, check.Equals, "POST")
	c.Assert(req.URL.Path, check.Equals, "/jobs")
	c.Assert(req.Header.Get("Content-Type"), check.Equals, "application/xml")
	var reqJob Job

	err = xml.Unmarshal(data, &reqJob)

	c.Assert(err, check.IsNil)
	c.Assert(reqJob, check.DeepEquals, myJob)

}

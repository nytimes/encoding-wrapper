package encodingdotcom

import (
	"encoding/json"
	"io/ioutil"
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

func (s *S) TestYesNoBoolean(c *check.C) {
	bTrue := YesNoBoolean(true)
	bFalse := YesNoBoolean(false)
	data, err := json.Marshal(bTrue)
	c.Assert(err, check.IsNil)
	c.Assert(string(data), check.Equals, `"yes"`)
	data, err = json.Marshal(bFalse)
	c.Assert(err, check.IsNil)
	c.Assert(string(data), check.Equals, `"no"`)
}

func (s *S) TestDoRequirementsParameters(c *check.C) {
	var req *http.Request
	var data string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
		data = req.FormValue("json")
		w.Write([]byte("it worked"))
	}))
	defer server.Close()
	client := Client{Endpoint: server.URL}
	resp, err := client.do(&Request{
		UserID:  "myuser",
		UserKey: "123",
		Action:  "AddMedia",
		MediaID: "123456",
		Source:  []string{"http://some.non.existent/video.mp4"},
	})
	c.Assert(err, check.IsNil)
	c.Assert(req, check.NotNil)
	c.Assert(req.Method, check.Equals, "POST")
	c.Assert(req.URL.Path, check.Equals, "/")
	c.Assert(req.Header.Get("Content-Type"), check.Equals, "application/x-www-form-urlencoded")
	var m map[string]interface{}
	err = json.Unmarshal([]byte(data), &m)
	c.Assert(err, check.IsNil)
	c.Assert(m, check.DeepEquals, map[string]interface{}{
		"query": map[string]interface{}{
			"userid":  "myuser",
			"userkey": "123",
			"action":  "AddMedia",
			"mediaid": "123456",
			"source":  []interface{}{"http://some.non.existent/video.mp4"},
		},
	})
	c.Assert(resp.StatusCode, check.Equals, http.StatusOK)
	respData, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	c.Assert(string(respData), check.Equals, "it worked")
}

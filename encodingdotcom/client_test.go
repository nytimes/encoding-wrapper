package encodingdotcom

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"gopkg.in/check.v1"
)

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
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.do(&Request{
		UserID:  client.UserID,
		UserKey: client.UserKey,
		Action:  "GetStatus",
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
			"action":  "GetStatus",
		},
	})
	c.Assert(resp.StatusCode, check.Equals, http.StatusOK)
	respData, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	c.Assert(string(respData), check.Equals, "it worked")
}

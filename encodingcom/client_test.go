package encodingcom

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"gopkg.in/check.v1"
)

func (s *S) mockErrorResponseObject(message *string, errors []string) interface{} {
	errorResponse := map[string]interface{}{
		"response": map[string]interface{}{},
	}
	if message != nil {
		errorResponse["response"].(map[string]interface{})["message"] = *message
	}
	if len(errors) > 0 {
		errorObject := make(map[string][]string)
		errorObject["error"] = errors
		errorResponse["response"].(map[string]interface{})["errors"] = errorObject
	}
	return errorResponse
}

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

func (s *S) TestDoMissingRequiredParameters(c *check.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		byteResponse, _ := json.Marshal(s.mockErrorResponseObject(nil, []string{"Wrong user id or key!"}))
		w.Write(byteResponse)
	}))
	defer server.Close()
	client := Client{Endpoint: server.URL}
	err := client.do(&request{
		Action:  "AddMedia",
		MediaID: "123456",
		Source:  []string{"http://some.non.existent/video.mp4"},
	}, nil)
	c.Assert(err, check.NotNil)
	apiErr, ok := err.(*APIError)
	c.Assert(ok, check.Equals, true)
	c.Assert(apiErr.Message, check.Equals, "")
	c.Assert(apiErr.Errors, check.DeepEquals, []string{"Wrong user id or key!"})
}

func (s *S) TestDoRequiredParameters(c *check.C) {
	var req *http.Request
	var data string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
		data = req.FormValue("json")
		w.Write([]byte(`{"response": {"status": "added"}}`))
	}))
	defer server.Close()
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	var respObj map[string]interface{}
	err := client.do(&request{
		UserID:  client.UserID,
		UserKey: client.UserKey,
		Action:  "GetStatus",
	}, &respObj)
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
	c.Assert(respObj, check.DeepEquals, map[string]interface{}{
		"response": map[string]interface{}{
			"status": "added",
		},
	})
}

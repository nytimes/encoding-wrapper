package encodingcom

import (
	"encoding/json"
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

func (s *S) TestDoGenericAction(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "Deleted"}}`, http.StatusOK)
	defer server.Close()
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	cancelMediaResponse, err := client.doGenericAction("12345", "CancelMedia")
	c.Assert(err, check.IsNil)
	c.Assert(cancelMediaResponse, check.DeepEquals, &Response{
		Message: "Deleted",
	})
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "CancelMedia")
}

func (s *S) TestDoMissingRequiredParameters(c *check.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		byteResponse, _ := json.Marshal(s.mockGenericResponseObject("", []string{"Wrong user id or key!"}))
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

func (s *S) TestDoGenericResponse(c *check.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		byteResponse, _ := json.Marshal(s.mockGenericResponseObject("it worked!", nil))
		w.Write(byteResponse)
	}))
	defer server.Close()
	client := Client{Endpoint: server.URL}
	var result map[string]*Response
	err := client.do(&request{
		Action:  "GetStatus",
		MediaID: "123456",
	}, &result)
	c.Assert(err, check.IsNil)
	c.Assert(result["response"].Message, check.Equals, "it worked!")
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

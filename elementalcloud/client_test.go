package elementalcloud

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

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
		Host:           "https://mycluster.cloud.elementaltechnologies.com",
		UserID:         "myuser",
		APIKey:         "secret-key",
		ExpirationTime: 45,
	}
	got := NewClient("https://mycluster.cloud.elementaltechnologies.com", "myuser", "secret-key", 45)
	c.Assert(*got, check.DeepEquals, expected)
}

func (s *S) TestCreateAuthKey(c *check.C) {
	path := "/jobs"
	userID := "myuser"
	APIKey := "api-key"
	expire := time.Unix(1, 0)
	innerKeyMD5 := md5.Sum([]byte(path + userID + APIKey + string(expire.Unix())))
	innerKey2 := hex.EncodeToString(innerKeyMD5[:])
	value := md5.Sum([]byte(APIKey + innerKey2))
	expected := hex.EncodeToString(value[:])
	client := NewClient("https://mycluster.cloud.elementaltechnologies.com", userID, APIKey, 45)
	got := client.createAuthKey(path, expire)
	c.Assert(got, check.Equals, expected)
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
	client := NewClient(server.URL, "myuser", "secret-key", 45)

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

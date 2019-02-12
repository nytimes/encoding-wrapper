package elementalconductor

import (
	// #nosec
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	expected := Client{
		Host:            "https://mycluster.cloud.elementaltechnologies.com",
		UserLogin:       "myuser",
		APIKey:          "elemental-secret-key",
		AuthExpires:     45,
		AccessKeyID:     "aws-access-key",
		SecretAccessKey: "aws-secret-key",
		Destination:     "destination",
	}
	got := NewClient("https://mycluster.cloud.elementaltechnologies.com", "myuser", "elemental-secret-key", 45, "aws-access-key", "aws-secret-key", "destination")
	if !reflect.DeepEqual(*got, expected) {
		t.Errorf("wrong client returned\nwant %#v\ngot  %#v", expected, *got)
	}
}

func TestCreateAuthKey(t *testing.T) {
	path := "/jobs"
	userID := "myuser"
	APIKey := "api-key"
	expire := time.Unix(1, 0)
	expireTimestamp := getUnixTimestamp(expire)
	// #nosec
	innerKeyMD5 := md5.Sum([]byte(path + userID + APIKey + expireTimestamp))
	innerKey2 := hex.EncodeToString(innerKeyMD5[:])
	// #nosec
	value := md5.Sum([]byte(APIKey + innerKey2))
	expected := hex.EncodeToString(value[:])
	client := NewClient("https://mycluster.cloud.elementaltechnologies.com", userID, APIKey, 45, "aws-access-key", "aws-secret-key", "destination")
	got := client.createAuthKey(path, expire)
	if got != expected {
		t.Errorf("wrong auth key returned\nwant %q\ngot  %q", expected, got)
	}
}

func TestDoRequiredParameters(t *testing.T) {
	var req *http.Request
	var data []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
		data, _ = ioutil.ReadAll(r.Body)
		w.Write([]byte(`<response>test</response>`))
	}))
	defer server.Close()
	client := NewClient(server.URL, "myuser", "elemental-secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	var respObj interface{}
	myJob := Job{
		XMLName: xml.Name{
			Local: "job",
		},
		Input: Input{
			FileInput: Location{
				URI:      "http://another.non.existent/video.mp4",
				Username: "user",
				Password: "pass123",
			},
		},
	}
	err := client.do("POST", "/jobs", myJob, &respObj)

	if err != nil {
		t.Fatal(err)
	}

	const (
		expectedMethod = "POST"
		expectedPath   = "/api/jobs"
	)
	expectedHeaders := map[string]string{
		"Accept":       "application/xml",
		"Content-type": "application/xml",
		"X-Auth-User":  client.UserLogin,
	}
	if req.Method != expectedMethod {
		t.Errorf("wrong http method used\nwant %q\ngot  %q", expectedMethod, req.Method)
	}
	if req.URL.Path != expectedPath {
		t.Errorf("wrong request path\nwant %q\ngot  %q", expectedPath, req.URL.Path)
	}

	for k, expectedValue := range expectedHeaders {
		v := req.Header.Get(k)
		if v != expectedValue {
			t.Errorf("wrong value for the header %q\nwant %q\ngot  %q", k, expectedValue, v)
		}
	}

	timestampInt, err := strconv.ParseInt(req.Header.Get("X-Auth-Expires"), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	timestampTime := time.Unix(timestampInt, 0)

	expectedAuthKey := client.createAuthKey("/jobs", timestampTime)
	if authKey := req.Header.Get("X-Auth-Key"); authKey != expectedAuthKey {
		t.Errorf("wrong auth key\nwant %q\ngot  %q", expectedAuthKey, authKey)
	}

	var reqJob Job
	err = xml.Unmarshal(data, &reqJob)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(reqJob, myJob) {
		t.Errorf("wrong job in the payload\nwant %#v\ngot  %#v", myJob, reqJob)
	}
}

func TestInvalidAuth(t *testing.T) {
	errorResponse := `<?xml version="1.0" encoding="UTF-8"?>
<errors>
  <error>You must be logged in to access this page.</error>
</errors>`
	server, _ := startServer(http.StatusUnauthorized, errorResponse)
	defer server.Close()
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	getJobsResponse, err := client.GetJob("1")
	if getJobsResponse != nil {
		t.Errorf("unexpected non-nil jobs response: %#v", getJobsResponse)
	}

	expectedAPIErr := &APIError{
		Status: http.StatusUnauthorized,
		Errors: errorResponse,
	}
	apiErr := err.(*APIError)
	if !reflect.DeepEqual(apiErr, expectedAPIErr) {
		t.Errorf("wrong api error returned\nwant %#v\ngot  %#v", expectedAPIErr, apiErr)
	}
}

func TestAPIErrorMarshalling(t *testing.T) {
	err := &APIError{
		Status: http.StatusInternalServerError,
		Errors: "something went wrong",
	}
	expectedError := `Error returned by the Elemental Conductor REST Interface: {"status":500,"errors":"something went wrong"}`
	if err.Error() != expectedError {
		t.Errorf("wrong error message\nwant %q\ngot  %q", expectedError, err.Error())
	}
}

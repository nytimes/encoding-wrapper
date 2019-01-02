package encodingcom

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	expected := &Client{
		Endpoint: "https://manage.encoding.com",
		UserID:   "myuser",
		UserKey:  "secret-key",
	}
	got, err := NewClient("https://manage.encoding.com", "myuser", "secret-key")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("wrong client returned\nwant %#v\ngot  %#v", expected, got)
	}
}

func TestYesNoBooleanMarshal(t *testing.T) {
	var tests = []struct {
		name     string
		input    bool
		expected string
	}{
		{
			"true",
			true,
			`"yes"`,
		},
		{
			"false",
			false,
			`"no"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := YesNoBoolean(test.input)
			data, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			r := string(data)
			if r != test.expected {
				t.Errorf("wrong result returned\nwant %v\ngot  %v", test.expected, r)
			}
		})
	}
}

func TestYesNoBooleanUnmarshal(t *testing.T) {
	data := []byte(`{"true":"yes", "false":"no"}`)
	var m map[string]YesNoBoolean
	err := json.Unmarshal(data, &m)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]YesNoBoolean{
		"true":  YesNoBoolean(true),
		"false": YesNoBoolean(false),
	}
	if !reflect.DeepEqual(m, expected) {
		t.Errorf("wrong value returned\nwant %#v\ngot  %#v", expected, m)
	}

	invalidData := []byte(`{"true":"true"}`)
	err = json.Unmarshal(invalidData, &m)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	expectedErrMsg := `invalid value: "true"`
	if err.Error() != expectedErrMsg {
		t.Errorf("wrong error message returned\nwant %q\ngot  %q", expectedErrMsg, err.Error())
	}
}

func TestZeroOneBooleanMarshal(t *testing.T) {
	var tests = []struct {
		name     string
		input    bool
		expected string
	}{
		{
			"true",
			true,
			`"1"`,
		},
		{
			"false",
			false,
			`"0"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := ZeroOneBoolean(test.input)
			data, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			r := string(data)
			if r != test.expected {
				t.Errorf("wrong result returned\nwant %v\ngot  %v", test.expected, r)
			}
		})
	}
}

func TestZeroOneBooleanUnmarshal(t *testing.T) {
	data := []byte(`{"true":"1", "false":"0"}`)
	var m map[string]ZeroOneBoolean
	err := json.Unmarshal(data, &m)
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]ZeroOneBoolean{
		"true":  ZeroOneBoolean(true),
		"false": ZeroOneBoolean(false),
	}
	if !reflect.DeepEqual(m, expected) {
		t.Errorf("wrong data returned\nwant %#v\ngot  %#v", expected, m)
	}

	invalidData := []byte(`{"true":"true"}`)
	err = json.Unmarshal(invalidData, &m)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	expectedErrMsg := `invalid value: "true"`
	if err.Error() != expectedErrMsg {
		t.Errorf("wrong error message returned\nwant %q\ngot  %q", expectedErrMsg, err.Error())
	}
}

func TestDoMediaAction(t *testing.T) {
	server, requests := startServer(`{"response": {"message": "Deleted"}}`)
	defer server.Close()
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	cancelMediaResponse, err := client.doMediaAction("12345", "CancelMedia")
	if err != nil {
		t.Fatal(err)
	}
	expectedResponse := &Response{Message: "Deleted"}
	if !reflect.DeepEqual(cancelMediaResponse, expectedResponse) {
		t.Errorf("wrong response returned\nwant %#v\ngot  %#v", expectedResponse, cancelMediaResponse)
	}
	req := <-requests
	const expectedAction = "CancelMedia"
	if req.query["action"] != expectedAction {
		t.Errorf("wrong action sent\nwant %q\ngot  %q", expectedAction, req.query["action"])
	}
}

func TestDoMediaActionFailure(t *testing.T) {
	server, requests := startServer(`{"response": {"message": "Deleted", "errors": {"error": "something went wrong"}}}`)
	defer server.Close()
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	cancelMediaResponse, err := client.doMediaAction("12345", "CancelMedia")
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if cancelMediaResponse != nil {
		t.Errorf("unexpected non-nil media response: %#v", cancelMediaResponse)
	}
	req := <-requests
	const expectedAction = "CancelMedia"
	if req.query["action"] != expectedAction {
		t.Errorf("wrong action sent\nwant %q\ngot  %q", expectedAction, req.query["action"])
	}
}

func TestDoMissingRequiredParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		byteResponse, _ := json.Marshal(mockMediaResponseObject("", "Wrong user id or key!"))
		w.Write(byteResponse)
	}))
	defer server.Close()
	client := Client{Endpoint: server.URL}
	err := client.do(&request{
		Action:  "AddMedia",
		MediaID: "123456",
		Source:  []string{"http://some.non.existent/video.mp4"},
	}, nil)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	expectedAPIErr := APIError{
		Message: "",
		Errors:  []string{"Wrong user id or key!"},
	}
	apiErr := err.(*APIError)
	if !reflect.DeepEqual(*apiErr, expectedAPIErr) {
		t.Errorf("wrong api error\nwant %#v\ngot  %#v", expectedAPIErr, *apiErr)
	}
}

func TestDoMediaResponse(t *testing.T) {
	const msg = "it worked!"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		byteResponse, _ := json.Marshal(mockMediaResponseObject(msg, ""))
		w.Write(byteResponse)
	}))
	defer server.Close()
	client := Client{Endpoint: server.URL}
	var result map[string]*Response
	err := client.do(&request{
		Action:  "GetStatus",
		MediaID: "123456",
	}, &result)
	if err != nil {
		t.Fatal(err)
	}
	if gotMsg := result["response"].Message; gotMsg != msg {
		t.Errorf("wrong message\nwant %q\ngot  %q", msg, gotMsg)
	}
}

func TestDoRequiredParameters(t *testing.T) {
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
	err := client.do(&request{Action: "GetStatus"}, &respObj)
	if err != nil {
		t.Fatal(err)
	}
	if req == nil {
		t.Fatal("unexpected <nil> request")
	}

	const (
		expectedMethod      = "POST"
		expectedPath        = "/"
		expectedContentType = "application/x-www-form-urlencoded"
	)
	if req.Method != expectedMethod {
		t.Errorf("wrong request method\nwant %q\ngot  %q", expectedMethod, req.Method)
	}
	if req.URL.Path != expectedPath {
		t.Errorf("wrong path\nwant %q\ngot  %q", expectedPath, req.URL.Path)
	}
	if ct := req.Header.Get("Content-Type"); ct != expectedContentType {
		t.Errorf("wrong Content-Type\nwant %q\ngot  %q", expectedContentType, ct)
	}

	var m map[string]interface{}
	err = json.Unmarshal([]byte(data), &m)
	if err != nil {
		t.Fatal(err)
	}
	expectedPayload := map[string]interface{}{
		"query": map[string]interface{}{
			"userid":  "myuser",
			"userkey": "123",
			"action":  "GetStatus",
		},
	}
	if !reflect.DeepEqual(m, expectedPayload) {
		t.Errorf("wrong payload sent\nwant %#v\ngot  %#v", expectedPayload, m)
	}
	expectedResponse := map[string]interface{}{
		"response": map[string]interface{}{
			"status": "added",
		},
	}
	if !reflect.DeepEqual(respObj, expectedResponse) {
		t.Errorf("wrong response obj\nwant %#v\ngot  %#v", expectedResponse, respObj)
	}
}

func TestDoInvalidResponse(t *testing.T) {
	server, _ := startServer(`{invalid json}`)
	defer server.Close()
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	var resp Response
	err := client.do(&request{Action: "GetStatus"}, &resp)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestAPIErrorRepresentation(t *testing.T) {
	err := &APIError{
		Message: "something went wrong",
		Errors:  []string{"error 1", "error 2"},
	}
	expectedMsg := `Error returned by the Encoding.com API: {"Message":"something went wrong","Errors":["error 1","error 2"]}`
	if err.Error() != expectedMsg {
		t.Errorf("wrong error message\nwant %q\ngot  %q", expectedMsg, err.Error())
	}
}

func mockMediaResponseObject(message string, errors string) interface{} {
	return map[string]interface{}{
		"response": map[string]interface{}{
			"message": message,
			"errors":  map[string]string{"error": errors},
		},
	}
}

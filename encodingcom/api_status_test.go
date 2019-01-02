package encodingcom

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAPIStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"Encoding Queue Processing Delays","status_code":"queue_slow","incident":"Our encoding queue is processing slower than normal.  Check back for updates."}`))
	}))
	defer server.Close()
	resp, err := APIStatus(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	expectedResp := APIStatusResponse{
		Status:     "Encoding Queue Processing Delays",
		StatusCode: "queue_slow",
		Incident:   "Our encoding queue is processing slower than normal.  Check back for updates.",
	}
	if !reflect.DeepEqual(*resp, expectedResp) {
		t.Errorf("wrong response returned\nwant %#v\ngot  %#v", expectedResp, *resp)
	}
}

func TestAPIStatusFailToConnect(t *testing.T) {
	_, err := APIStatus("http://192.0.2.13:8080")
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestAPIStatusInvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{not a valid json}`))
	}))
	defer server.Close()
	_, err := APIStatus(server.URL)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestAPIStatusOK(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"ok", true},
		{"encoding_delay", false},
		{"api_out", false},
		{"maintenance", false},
		{"pc_queue_slow", false},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			status := APIStatusResponse{StatusCode: test.input}
			got := status.OK()
			if got != test.want {
				t.Errorf("wrong status returned\nwant %v\ngot  %v", test.want, got)
			}
		})
	}
}

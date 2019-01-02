package encodingcom

import (
	"reflect"
	"testing"
	"time"
)

func TestGetStatusSingle(t *testing.T) {
	server, requests := startServer(`
{
	"response": {
        "job": {
				"id": "abc123",
				"userid": "myuser",
				"sourcefile": "http://some.video/file.mp4",
				"status": "Finished",
				"notifyurl": "http://ping.me/please",
				"created": "2015-12-31 20:45:30",
				"started": "2015-12-31 20:45:34",
				"finished": "2015-12-31 21:00:03",
				"prevstatus": "Saving",
				"downloaded": "2015-12-31 20:45:32",
				"uploaded": "2015-12-31 20:59:54",
				"time_left": "0",
				"progress": "100",
				"time_left_current": "0",
				"progress_current": "100.0",
				"format": {
						"id": "f123",
						"status": "Finished",
						"description": "Something",
						"created": "2015-12-31 20:45:30",
						"started": "2015-12-31 20:45:34",
						"finished": "2015-12-31 21:00:03",
						"s3_destination": "https://s3.amazonaws.com/not-really/valid.mp4",
						"cf_destination": "https://blablabla.cloudfront.net/not-valid.mp4",
						"convertedsize": "65723",
						"destination": "s3://mynicebucket",
						"destination_status": "Saved"
					}
			}
	}
}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	status, err := client.GetStatus([]string{"abc123"}, true)
	if err != nil {
		t.Fatal(err)
	}

	expectedCreateDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:30")
	expectedStartDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:34")
	expectedFinishDate, _ := time.Parse(dateTimeLayout, "2015-12-31 21:00:03")
	expectedDownloadDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:32")
	expectedUploadDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:59:54")
	expected := []StatusResponse{
		{
			MediaID:             "abc123",
			UserID:              "myuser",
			SourceFile:          "http://some.video/file.mp4",
			MediaStatus:         "Finished",
			PreviousMediaStatus: "Saving",
			NotifyURL:           "http://ping.me/please",
			CreateDate:          expectedCreateDate,
			StartDate:           expectedStartDate,
			FinishDate:          expectedFinishDate,
			DownloadDate:        expectedDownloadDate,
			UploadDate:          expectedUploadDate,
			TimeLeft:            "0",
			Progress:            100.0,
			TimeLeftCurrentJob:  "0",
			ProgressCurrentJob:  100.0,
			Formats: []FormatStatus{
				{
					ID:            "f123",
					Status:        "Finished",
					Description:   "Something",
					CreateDate:    expectedCreateDate,
					StartDate:     expectedStartDate,
					FinishDate:    expectedFinishDate,
					S3Destination: "https://s3.amazonaws.com/not-really/valid.mp4",
					CFDestination: "https://blablabla.cloudfront.net/not-valid.mp4",
					FileSize:      "65723",
					Destinations:  []DestinationStatus{{Name: "s3://mynicebucket", Status: "Saved"}},
				},
			},
		},
	}
	if !reflect.DeepEqual(status, expected) {
		t.Errorf("wrong status returned\nwant %#v\ngot  %#v", expected, status)
	}

	req := <-requests
	expectedQuery := map[string]interface{}{
		"action":   "GetStatus",
		"mediaid":  "abc123",
		"extended": "yes",
		"userid":   "myuser",
		"userkey":  "123",
	}
	if !reflect.DeepEqual(req.query, expectedQuery) {
		t.Errorf("wrong query\nwant %#v\ngot  %#v", expectedQuery, req.query)
	}
}

func TestGetStatusMultiple(t *testing.T) {
	server, requests := startServer(`
{
	"response": {
		"job": [
			{
				"id": "abc123",
				"userid": "myuser",
				"sourcefile": "http://some.video/file.mp4",
				"status": "Finished",
				"notifyurl": "http://ping.me/please",
				"created": "2015-12-31 20:45:30",
				"started": "2015-12-31 20:45:34",
				"finished": "2015-12-31 21:00:03",
				"prevstatus": "Saving",
				"downloaded": "2015-12-31 20:45:32",
				"uploaded": "2015-12-31 20:59:54",
				"time_left": "0",
				"progress": "100",
				"time_left_current": "0",
				"progress_current": "100.0",
				"format": [
					{
						"id": "f123",
						"status": "Finished",
						"created": "2015-12-31 20:45:30",
						"started": "2015-12-31 20:45:34",
						"finished": "2015-12-31 21:00:03",
						"s3_destination": "https://s3.amazonaws.com/not-really/valid.mp4",
						"cf_destination": "https://blablabla.cloudfront.net/not-valid.mp4",
						"convertedsize": "65724",
						"destination": [
							null,
							"s3://myunclebucket/file.mp4"
						],
						"destination_status": [
							null,
							"Saved"
						]
					},
					{
						"id": "f124",
						"status": "Finished",
						"created": "2015-12-31 20:45:30",
						"started": "2015-12-31 20:45:34",
						"finished": "2015-12-31 21:00:03",
						"s3_destination": "https://s3.amazonaws.com/not-really/valid.mp4",
						"cf_destination": "https://blablabla.cloudfront.net/not-valid.mp4",
						"convertedsize": "65725",
						"destination": [
							"s3://mynicebucket/file.mp4",
							"s3://myunclebucket/file.mp4"
						],
						"destination_status": null
					}
				]
			},
			{
				"id": "abc124",
				"userid": "myuser",
				"sourcefile": "http://some.video/file.mp4",
				"status": "Finished",
				"notifyurl": "http://ping.me/please",
				"created": "2015-12-31 20:45:30",
				"started": "2015-12-31 20:45:34",
				"finished": "2015-12-31 21:00:03",
				"prevstatus": "Saving",
				"downloaded": "2015-12-31 20:45:32",
				"uploaded": "2015-12-31 20:59:54",
				"time_left": "0",
				"progress": "100",
				"time_left_current": "0",
				"progress_current": "100.0",
				"format": {
						"id": "f123",
						"status": "Finished",
						"created": "2015-12-31 20:45:30",
						"started": "2015-12-31 20:45:34",
						"finished": "2015-12-31 21:00:03",
						"s3_destination": "https://s3.amazonaws.com/not-really/valid.mp4",
						"cf_destination": "https://blablabla.cloudfront.net/not-valid.mp4",
						"convertedsize": "65726",
						"destination": [
							"s3://mynicebucket/file.mp4",
							"s3://myunclebucket/file.mp4"
						],
						"destination_status": [
							"Saved",
							"Saved"
						]
					}
			}
		]
	}
}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	status, err := client.GetStatus([]string{"abc123", "abc124"}, true)
	if err != nil {
		t.Fatal(err)
	}

	expectedCreateDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:30")
	expectedStartDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:34")
	expectedFinishDate, _ := time.Parse(dateTimeLayout, "2015-12-31 21:00:03")
	expectedDownloadDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:32")
	expectedUploadDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:59:54")
	expected := []StatusResponse{
		{
			MediaID:             "abc123",
			UserID:              "myuser",
			SourceFile:          "http://some.video/file.mp4",
			MediaStatus:         "Finished",
			PreviousMediaStatus: "Saving",
			NotifyURL:           "http://ping.me/please",
			CreateDate:          expectedCreateDate,
			StartDate:           expectedStartDate,
			FinishDate:          expectedFinishDate,
			DownloadDate:        expectedDownloadDate,
			UploadDate:          expectedUploadDate,
			TimeLeft:            "0",
			Progress:            100.0,
			TimeLeftCurrentJob:  "0",
			ProgressCurrentJob:  100.0,
			Formats: []FormatStatus{
				{
					ID:            "f123",
					Status:        "Finished",
					CreateDate:    expectedCreateDate,
					StartDate:     expectedStartDate,
					FinishDate:    expectedFinishDate,
					S3Destination: "https://s3.amazonaws.com/not-really/valid.mp4",
					CFDestination: "https://blablabla.cloudfront.net/not-valid.mp4",
					FileSize:      "65724",
					Destinations: []DestinationStatus{
						{Name: "", Status: ""},
						{Name: "s3://myunclebucket/file.mp4", Status: "Saved"},
					},
				},
				{
					ID:            "f124",
					Status:        "Finished",
					CreateDate:    expectedCreateDate,
					StartDate:     expectedStartDate,
					FinishDate:    expectedFinishDate,
					S3Destination: "https://s3.amazonaws.com/not-really/valid.mp4",
					CFDestination: "https://blablabla.cloudfront.net/not-valid.mp4",
					FileSize:      "65725",
					Destinations: []DestinationStatus{
						{Name: "s3://mynicebucket/file.mp4", Status: ""},
						{Name: "s3://myunclebucket/file.mp4", Status: ""},
					},
				},
			},
		},
		{
			MediaID:             "abc124",
			UserID:              "myuser",
			SourceFile:          "http://some.video/file.mp4",
			MediaStatus:         "Finished",
			PreviousMediaStatus: "Saving",
			NotifyURL:           "http://ping.me/please",
			CreateDate:          expectedCreateDate,
			StartDate:           expectedStartDate,
			FinishDate:          expectedFinishDate,
			DownloadDate:        expectedDownloadDate,
			UploadDate:          expectedUploadDate,
			TimeLeft:            "0",
			Progress:            100.0,
			TimeLeftCurrentJob:  "0",
			ProgressCurrentJob:  100.0,
			Formats: []FormatStatus{
				{
					ID:            "f123",
					Status:        "Finished",
					CreateDate:    expectedCreateDate,
					StartDate:     expectedStartDate,
					FinishDate:    expectedFinishDate,
					S3Destination: "https://s3.amazonaws.com/not-really/valid.mp4",
					CFDestination: "https://blablabla.cloudfront.net/not-valid.mp4",
					FileSize:      "65726",
					Destinations: []DestinationStatus{
						{Name: "s3://mynicebucket/file.mp4", Status: "Saved"},
						{Name: "s3://myunclebucket/file.mp4", Status: "Saved"},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(status, expected) {
		t.Errorf("wrong status returned\nwant %#v\ngot  %#v", expected, status)
	}

	req := <-requests
	expectedQuery := map[string]interface{}{
		"action":   "GetStatus",
		"mediaid":  "abc123,abc124",
		"extended": "yes",
		"userid":   "myuser",
		"userkey":  "123",
	}
	if !reflect.DeepEqual(req.query, expectedQuery) {
		t.Errorf("wrong query\nwant %#v\ngot  %#v", expectedQuery, req.query)
	}
}

// Some data are only available when extended=no
func TestGetStatusNotExtended(t *testing.T) {
	server, requests := startServer(nonExtendedStatus)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	status, err := client.GetStatus([]string{"abc123"}, false)
	if err != nil {
		t.Fatal(err)
	}

	expectedCreateDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:30")
	expectedStartDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:34")
	expectedFinishDate, _ := time.Parse(dateTimeLayout, "2015-12-31 21:00:03")
	expectedDownloadDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:32")
	expected := []StatusResponse{
		{
			MediaID:      "abc123",
			UserID:       "myuser",
			SourceFile:   "http://some.video/file.mp4",
			MediaStatus:  "Finished",
			CreateDate:   expectedCreateDate,
			StartDate:    expectedStartDate,
			FinishDate:   expectedFinishDate,
			DownloadDate: expectedDownloadDate,
			TimeLeft:     "21",
			Progress:     100.0,
			Formats: []FormatStatus{
				{
					ID:           "f123",
					Status:       "Finished",
					CreateDate:   expectedCreateDate,
					StartDate:    expectedStartDate,
					FinishDate:   expectedFinishDate,
					Destinations: []DestinationStatus{{Name: "s3://mynicebucket", Status: "Saved"}},
					Size:         "0x1080",
					Bitrate:      "3500k",
					Output:       "mp4",
					VideoCodec:   "libx264",
					AudioCodec:   "dolby_aac",
					FileSize:     "78544430",
				},
			},
		},
	}
	if !reflect.DeepEqual(status, expected) {
		t.Errorf("wrong status returned\nwant %#v\ngot  %#v", expected, status)
	}

	req := <-requests
	expectedQuery := map[string]interface{}{
		"action":  "GetStatus",
		"mediaid": "abc123",
		"userid":  "myuser",
		"userkey": "123",
	}
	if !reflect.DeepEqual(req.query, expectedQuery) {
		t.Errorf("wrong query\nwant %#v\ngot  %#v", expectedQuery, req.query)
	}
}

func TestGetStatusZeroTime(t *testing.T) {
	server, _ := startServer(`
{
	"response": {
		"job": {
			"id":"abc123",
			"userid":"myuser",
			"sourcefile":"http://some.file/wait-wat",
			"status":"Error",
			"created":"2016-01-29 19:32:32",
			"started":"2016-01-29 19:32:32",
			"finished":"0000-00-00 00:00:00",
			"downloaded":"0000-00-00 00:00:00",
			"description":"Download error:  The requested URL returned error: 403 Forbidden",
			"processor":"AMAZON",
			"region":"oak-private-clive",
			"time_left":"50",
			"progress":"50.0",
			"time_left_current":"0",
			"progress_current":"0.0",
			"format":{
				"id":"164478401",
				"status":"New",
				"created":"2016-01-29 19:32:32",
				"started":"0000-00-00 00:00:00",
				"finished":"0000-00-00 00:00:00",
				"destination":"http://s4.amazonaws.com/future",
				"destination_status":"Open",
				"convertedsize":"0",
				"queued":"0000-00-00 00:00:00",
				"converttime":"0",
				"time_left":"40",
				"progress":"0.0",
				"time_left_current":"0",
				"progress_current":"0.0"
			},
			"queue_time":"0"
		}
	}
}`)
	defer server.Close()

	const mediaID = "abc123"
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	status, err := client.GetStatus([]string{mediaID}, true)
	if err != nil {
		t.Fatal(err)
	}

	expectedCreateDate, _ := time.Parse(dateTimeLayout, "2016-01-29 19:32:32")
	if status[0].MediaID != mediaID {
		t.Errorf("wrong media id returned\nwant %q\ngot  %q", mediaID, status[0].MediaID)
	}
	if !reflect.DeepEqual(status[0].CreateDate, expectedCreateDate) {
		t.Errorf("wrong create date\nwant %s\ngot  %s", expectedCreateDate, status[0].CreateDate)
	}
	if !status[0].FinishDate.IsZero() {
		t.Errorf("unexpected non-zero finish date: %#v", status[0].FinishDate)
	}
	if !status[0].DownloadDate.IsZero() {
		t.Errorf("unexpected non-zero download date: %#v", status[0].DownloadDate)
	}
}

func TestGetStatusNoMedia(t *testing.T) {
	var client Client
	status, err := client.GetStatus(nil, true)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	const expectedErrMsg = "please provide at least one media id"
	if err.Error() != expectedErrMsg {
		t.Errorf("wrong error message\nwant %q\ngot  %q", expectedErrMsg, err.Error())
	}
	if len(status) != 0 {
		t.Errorf("unexpected non-empty status response: %#v", status)
	}
}

func TestGetStatusError(t *testing.T) {
	server, _ := startServer(`{"response": {"message": "", "errors": {"error": "wait what?"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.GetStatus([]string{"some-media"}, true)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if resp != nil {
		t.Errorf("unexpected non-nil response: %#v", resp)
	}
}

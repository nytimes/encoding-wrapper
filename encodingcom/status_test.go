package encodingcom

import (
	"time"

	"gopkg.in/check.v1"
)

func (s *S) TestGetStatusSingle(c *check.C) {
	server, requests := s.startServer(`
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
						"destination": [
							"s3://mynicebucket"
						],
						"destination_status": [
							"Saved"
						]
					}
				]
			}
		]
	}
}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	status, err := client.GetStatus([]string{"abc123"})
	c.Assert(err, check.IsNil)

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
					Destinations:  []DestinationStatus{{Name: "s3://mynicebucket", Status: "Saved"}},
				},
			},
		},
	}
	c.Assert(status, check.DeepEquals, expected)

	req := <-requests
	c.Assert(req.query["action"], check.Equals, "GetStatus")
	c.Assert(req.query["mediaid"], check.Equals, "abc123")
	c.Assert(req.query["extended"], check.Equals, "yes")
}

func (s *S) TestGetStatusMultiple(c *check.C) {
	server, requests := s.startServer(`
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
						"destination": [
							"s3://mynicebucket"
						],
						"destination_status": [
							"Saved"
						]
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
				"format": [
					{
						"id": "f123",
						"status": "Finished",
						"created": "2015-12-31 20:45:30",
						"started": "2015-12-31 20:45:34",
						"finished": "2015-12-31 21:00:03",
						"s3_destination": "https://s3.amazonaws.com/not-really/valid.mp4",
						"cf_destination": "https://blablabla.cloudfront.net/not-valid.mp4",
						"destination": [
							"s3://mynicebucket"
						],
						"destination_status": [
							"Saved"
						]
					}
				]
			}
		]
	}
}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	status, err := client.GetStatus([]string{"abc123", "abc124"})
	c.Assert(err, check.IsNil)

	expectedCreateDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:30")
	expectedStartDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:34")
	expectedFinishDate, _ := time.Parse(dateTimeLayout, "2015-12-31 21:00:03")
	expectedDownloadDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:32")
	expectedUploadDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:59:54")
	status1 := StatusResponse{
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
				Destinations:  []DestinationStatus{{Name: "s3://mynicebucket", Status: "Saved"}},
			},
		},
	}
	status2 := status1
	status2.MediaID = "abc124"
	expected := []StatusResponse{status1, status2}
	c.Assert(status, check.DeepEquals, expected)

	req := <-requests
	c.Assert(req.query["action"], check.Equals, "GetStatus")
	c.Assert(req.query["mediaid"], check.Equals, "abc123,abc124")
	c.Assert(req.query["extended"], check.Equals, "yes")
}

func (s *S) TestGetStatusNoMedia(c *check.C) {
	var client Client
	status, err := client.GetStatus(nil)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "please provide at least one media id")
	c.Assert(status, check.HasLen, 0)
}

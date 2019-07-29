package encodingcom

import (
	"reflect"
	"testing"
	"time"
)

func TestAddMedia(t *testing.T) {
	server, requests := startServer(`{"response": {"message": "Added", "MediaID": "1234567"}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	format := Format{
		Output:       []string{"mp4"},
		VideoCodec:   "x264",
		AudioCodec:   "aac",
		Bitrate:      "900k",
		AudioBitrate: "64k",
	}
	addMediaResponse, err := client.AddMedia([]string{"http://another.non.existent/video.mov"},
		[]Format{format}, "us-east-1")

	if err != nil {
		t.Fatal(err)
	}
	expectedResp := &AddMediaResponse{
		Message: "Added",
		MediaID: "1234567",
	}
	if !reflect.DeepEqual(addMediaResponse, expectedResp) {
		t.Errorf("wrong response\nwant %#v\ngot  %#v", expectedResp, addMediaResponse)
	}
	const expectedAction = "AddMedia"
	req := <-requests
	if action := req.query["action"]; action != expectedAction {
		t.Errorf("wrong action returned\nwant %q\ngot  %q", expectedAction, action)
	}
}

func TestAddMediaError(t *testing.T) {
	server, _ := startServer(`{"response": {"message": "Added", "errors": {"error": "something went wrong"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	format := Format{
		Output:       []string{"mp4"},
		VideoCodec:   "x264",
		AudioCodec:   "aac",
		Bitrate:      "900k",
		AudioBitrate: "64k",
	}
	addMediaResponse, err := client.AddMedia([]string{"http://another.non.existent/video.mov"},
		[]Format{format}, "us-east-1")
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if addMediaResponse != nil {
		t.Errorf("unexpected non-nil response: %#v", addMediaResponse)
	}
}

func TestStopMedia(t *testing.T) {
	server, requests := startServer(`{"response": {"message": "Stopped"}}`)
	defer server.Close()

	const mediaID = "some-media"
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.StopMedia(mediaID)
	if err != nil {
		t.Fatal(err)
	}
	expectedResp := &Response{Message: "Stopped"}
	if !reflect.DeepEqual(resp, expectedResp) {
		t.Errorf("wrong response returned\nwant %#v\ngot  %#v", expectedResp, resp)
	}
	const expectedAction = "StopMedia"
	req := <-requests
	if action := req.query["action"]; action != expectedAction {
		t.Errorf("wrong action sent\nwant %q\ngot  %q", expectedAction, action)
	}
	if req.query["mediaid"] != mediaID {
		t.Errorf("wrong media id sent\nwant %q\ngot  %q", mediaID, req.query["mediaid"])
	}
}

func TestStopMediaError(t *testing.T) {
	server, requests := startServer(`{"response": {"message": "failed", "errors": {"error": "something went wrong"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.StopMedia("some-media")
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if resp != nil {
		t.Errorf("unexpected non-nil response: %#v", resp)
	}
	const expectedAction = "StopMedia"
	req := <-requests
	if action := req.query["action"]; action != expectedAction {
		t.Errorf("wrong action sent\nwant %q\ngot  %q", expectedAction, action)
	}
}

func TestCancelMedia(t *testing.T) {
	server, requests := startServer(`{"response": {"message": "Canceled"}}`)
	defer server.Close()

	const mediaID = "some-media"
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.CancelMedia(mediaID)
	if err != nil {
		t.Fatal(err)
	}
	expectedResp := &Response{Message: "Canceled"}
	if !reflect.DeepEqual(resp, expectedResp) {
		t.Errorf("wrong response returned\nwant %#v\ngot  %#v", expectedResp, resp)
	}
	const expectedAction = "CancelMedia"
	req := <-requests
	if action := req.query["action"]; action != expectedAction {
		t.Errorf("wrong action sent\nwant %q\ngot  %q", expectedAction, action)
	}
	if req.query["mediaid"] != mediaID {
		t.Errorf("wrong media id sent\nwant %q\ngot  %q", mediaID, req.query["mediaid"])
	}
}

func TestCancelMediaError(t *testing.T) {
	server, _ := startServer(`{"response": {"message": "failed", "errors": {"error": "something went wrong"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.CancelMedia("some-media")
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if resp != nil {
		t.Errorf("unexpected non-nil error: %#v", resp)
	}
}

func TestRestartMedia(t *testing.T) {
	server, requests := startServer(`{"response": {"message": "Restarted"}}`)
	defer server.Close()

	const mediaID = "some-media"
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	expectedResp := &Response{Message: "Restarted"}
	resp, err := client.RestartMedia(mediaID, false)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(resp, expectedResp) {
		t.Errorf("wrong response\nwant %#v\ngot  %#v", expectedResp, resp)
	}

	const expectedAction = "RestartMedia"
	req := <-requests
	if action := req.query["action"]; action != expectedAction {
		t.Errorf("wrong action sent\nwant %q\ngot  %q", expectedAction, action)
	}
	if req.query["mediaid"] != mediaID {
		t.Errorf("wrong media id sent\nwant %q\ngot  %q", mediaID, req.query["mediaid"])
	}
}

func TestRestartMediaWithErrors(t *testing.T) {
	server, requests := startServer(`{"response": {"message": "Restarted"}}`)
	defer server.Close()

	const mediaID = "some-media"
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	expectedResp := &Response{Message: "Restarted"}
	resp, err := client.RestartMedia(mediaID, true)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(resp, expectedResp) {
		t.Errorf("wrong response\nwant %#v\ngot  %#v", expectedResp, resp)
	}

	const expectedAction = "RestartMediaErrors"
	req := <-requests
	if action := req.query["action"]; action != expectedAction {
		t.Errorf("wrong action sent\nwant %q\ngot  %q", expectedAction, action)
	}
	if req.query["mediaid"] != mediaID {
		t.Errorf("wrong media id sent\nwant %q\ngot  %q", mediaID, req.query["mediaid"])
	}
}

func TestRestartMediaError(t *testing.T) {
	server, _ := startServer(`{"response": {"message": "failed", "errors": {"error": "something went wrong"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.RestartMedia("some-media", false)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if resp != nil {
		t.Errorf("unexpected non-nil response: %#v", resp)
	}
}

func TestRestartMediaTask(t *testing.T) {
	server, requests := startServer(`{"response": {"message": "Task restarted"}}`)
	defer server.Close()

	const (
		mediaID = "some-media"
		taskID  = "some-task"
	)
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.RestartMediaTask(mediaID, taskID)
	if err != nil {
		t.Fatal(err)
	}
	expectedResp := &Response{Message: "Task restarted"}
	if !reflect.DeepEqual(resp, expectedResp) {
		t.Errorf("wrong response returned\nwant %#v\ngot  %#v", expectedResp, resp)
	}

	const expectedAction = "RestartMediaTask"
	req := <-requests
	if action := req.query["action"]; action != expectedAction {
		t.Errorf("wrong action sent\nwant %q\ngot  %q", expectedAction, action)
	}
	if req.query["mediaid"] != mediaID {
		t.Errorf("wrong media id sent\nwant %q\ngot  %q", mediaID, req.query["mediaid"])
	}
	if req.query["taskid"] != taskID {
		t.Errorf("wrong task id sent\nwant %q\ngot  %q", taskID, req.query["taskid"])
	}
}

func TestRestartMediaTaskError(t *testing.T) {
	server, _ := startServer(`{"response": {"message": "Failed to restart", "errors": {"error": "something went really bad"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.RestartMediaTask("some-media", "some-task")
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if resp != nil {
		t.Errorf("unexpected non-nil response: %#v", resp)
	}
}

func TestListMedia(t *testing.T) {
	server, requests := startServer(`
{
    "response":{
        "media":[
            {
                "mediafile":"http://another.non.existent/video.mp4",
                "mediaid":"1234567",
                "mediastatus":"Finished",
                "createdate":"2015-12-31 20:45:30",
                "startdate":"2015-12-31 20:45:50",
                "finishdate":"2015-12-31 20:48:54"
            }
        ]
    }
}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	listMediaResponse, err := client.ListMedia()
	if err != nil {
		t.Fatal(err)
	}

	mockCreateDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:30")
	mockStartDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:50")
	mockFinishDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:48:54")

	expectedResp := &ListMediaResponse{
		Media: []ListMediaResponseItem{
			{
				"http://another.non.existent/video.mp4",
				"1234567",
				"Finished",
				MediaDateTime{mockCreateDate},
				MediaDateTime{mockStartDate},
				MediaDateTime{mockFinishDate},
			},
		},
	}
	if !reflect.DeepEqual(listMediaResponse, expectedResp) {
		t.Errorf("wrong response returned\nwant %#v\ngot  %#v", expectedResp, listMediaResponse)
	}

	const expectedAction = "GetMediaList"
	req := <-requests
	if action := req.query["action"]; action != expectedAction {
		t.Errorf("wrong action\nwant %q\ngot  %q", expectedAction, action)
	}
}

func TestListMediaError(t *testing.T) {
	server, _ := startServer(`{"response": {"message": "", "errors": {"error": "can't list"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.ListMedia()
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if resp != nil {
		t.Errorf("unexpected non-nil response: %#v", resp)
	}
}

func TestGetMediaInfo(t *testing.T) {
	server, requests := startServer(`
{
	"response": {
		"bitrate": "1807k",
		"duration": "6464.83",
		"audio_bitrate": "128k",
		"video_codec": "mpeg4",
		"video_bitrate": "1679k",
		"frame_rate": "23.98",
		"size": "640x352",
		"pixel_aspect_ratio": "1:1",
		"display_aspect_ratio": "20:11",
		"audio_codec": "ac3",
		"audio_sample_rate": "48000",
		"audio_channels": "2",
		"rotation":"90"
	}
}`)
	defer server.Close()

	const mediaID = "m-123"
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	mediaInfo, err := client.GetMediaInfo(mediaID)
	if err != nil {
		t.Fatal(err)
	}

	expectedMediaInfo := &MediaInfo{
		Bitrate:            "1807k",
		Duration:           6464*time.Second + time.Duration(0.83*float64(time.Second)),
		VideoCodec:         "mpeg4",
		VideoBitrate:       "1679k",
		Framerate:          "23.98",
		Size:               "640x352",
		PixelAspectRatio:   "1:1",
		DisplayAspectRatio: "20:11",
		AudioCodec:         "ac3",
		AudioSampleRate:    uint(48000),
		AudioChannels:      "2",
		AudioBitrate:       "128k",
		Rotation:           90,
	}
	if !reflect.DeepEqual(mediaInfo, expectedMediaInfo) {
		t.Errorf("wrong MediaInfo returned\nwant %#v\ngot  %#v", expectedMediaInfo, mediaInfo)
	}

	const expectedAction = "GetMediaInfo"
	req := <-requests
	if action := req.query["action"]; action != expectedAction {
		t.Errorf("wrong action\nwant %q\ngot  %q", expectedAction, action)
	}
	if req.query["mediaid"] != mediaID {
		t.Errorf("wrong media id sent\nwant %q\ngot  %q", mediaID, req.query["mediaid"])
	}
}

func TestGetMediaInfoError(t *testing.T) {
	server, _ := startServer(`{"response": {"message": "", "errors": {"error": "wait what?"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.GetMediaInfo("some-media")
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if resp != nil {
		t.Errorf("unexpected non-nil response: %#v", resp)
	}
}

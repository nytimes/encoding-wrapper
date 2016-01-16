package encodingcom

import (
	"net/http"
	"time"

	"gopkg.in/check.v1"
)

func (s *S) TestAddMedia(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "Added", "MediaID": "1234567"}}`, http.StatusOK)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	addMediaResponse, err := client.AddMedia([]string{"http://another.non.existent/video.mov"}, &Format{
		Output:       []string{"http://another.non.existent/video.mp4"},
		VideoCodec:   "x264",
		AudioCodec:   "aac",
		Bitrate:      "900k",
		AudioBitrate: "64k",
	})

	c.Assert(err, check.IsNil)
	c.Assert(addMediaResponse, check.DeepEquals, &AddMediaResponse{
		Message: "Added",
		MediaID: "1234567",
	})
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "AddMedia")
}

func (s *S) TestListMedia(c *check.C) {
	server, requests := s.startServer(`
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
}`, http.StatusOK)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	listMediaResponse, err := client.ListMedia()
	c.Assert(err, check.IsNil)

	mockCreateDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:30")
	mockStartDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:50")
	mockFinishDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:48:54")

	c.Assert(listMediaResponse, check.DeepEquals, &ListMediaResponse{
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
	})
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "GetMediaList")
}

func (s *S) TestGetMediaInfo(c *check.C) {
	server, requests := s.startServer(`
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
		"audio_channels": "2"
	}
}`, http.StatusOK)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	mediaInfo, err := client.GetMediaInfo("m-123")
	c.Assert(err, check.IsNil)

	c.Assert(mediaInfo, check.DeepEquals, &MediaInfo{
		Bitrate:            "1807k",
		Duration:           6464.83,
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
	})

	req := <-requests
	c.Assert(req.query["action"], check.Equals, "GetMediaInfo")
	c.Assert(req.query["mediaid"], check.Equals, "m-123")
}

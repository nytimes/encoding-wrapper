package encodingcom

import (
	"time"

	"gopkg.in/check.v1"
)

func (s *S) TestAddMedia(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "Added", "MediaID": "1234567"}}`, 200)
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
}`, 200)
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

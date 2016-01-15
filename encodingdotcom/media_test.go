package encodingdotcom

import (
	"net/http"
	"net/http/httptest"

	"gopkg.in/check.v1"
)

func (s *S) TestAddMedia(c *check.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
				"response": {
				"message": "Added",
				"MediaID": "1234567"
				}
		}`))
	}))

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
}

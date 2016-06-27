package encodingcom

import "gopkg.in/check.v1"

func (s *S) TestGetPresetsList(c *check.C) {
	server, requests := s.startServer(`
{
	"response": {
		"user": [{
			"name":"webm_1080p",
			"type":"user",
			"output":"webm",
			"format":{
				"output":"webm",
				"audio_bitrate":"192k",
				"audio_sample_rate":"48000",
				"bitrate":"5000k",
				"framerate":"30",
				"keep_aspect_ratio":"no",
				"set_aspect_ratio":"1280 x 1080",
				"video_codec":"libvpx",
				"video_codec_parameters":"no",
				"size":"1920x1080"
			}
		},
		{
			"name":"webm_480p",
			"type":"user",
			"output":"webm",
			"format":{
				"output":"webm",
				"audio_bitrate":"192k",
				"audio_sample_rate":"48000",
				"bitrate":"2500k",
				"framerate":"30",
				"video_codec":"libvpx",
				"video_codec_parameters":"no",
				"size":"1280x480"
			}
		}],
		"ui": [
		{
			"name":"3GP 144p 256k",
			"type":"ui",
			"output":"3gp",
			"format":{
				"output":"3gp",
				"size":"176x144",
				"bitrate":"256k",
				"audio_bitrate":"12.2k",
				"audio_sample_rate":"8000",
				"audio_channels_number":"1",
				"framerate":"15",
				"keep_aspect_ratio":"yes",
				"video_codec":"h263",
				"profile":"baseline",
				"audio_codec":"libamr_nb",
				"two_pass":"no",
				"turbo":"no",
				"twin_turbo":"no",
				"keyframe":"90",
				"audio_volume":"100",
				"rotate":"def",
				"strip_chapters":"no"
			}
		},
		{
			"name":"3GP 288p 360k",
			"type":"ui",
			"output":"3gp",
			"format":{
				"output":"3gp",
				"size":"352x288",
				"bitrate":"360k",
				"audio_bitrate":"12.2k",
				"audio_sample_rate":"8000",
				"audio_channels_number":"1",
				"framerate":"24",
				"keep_aspect_ratio":"yes",
				"video_codec":"h263",
				"profile":"baseline",
				"audio_codec":"libamr_nb",
				"two_pass":"no",
				"turbo":"no",
				"twin_turbo":"no",
				"keyframe":"120",
				"audio_volume":"100",
				"rotate":"def",
				"strip_chapters":"no"
			}
		}]
	}
}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.ListPresets(AllPresets)
	c.Assert(err, check.IsNil)

	c.Assert(resp, check.DeepEquals, &ListPresetsResponse{
		UserPresets: []Preset{
			{
				Name:   "webm_1080p",
				Type:   "user",
				Output: "webm",
				Format: PresetFormat{
					Output:               "webm",
					AudioBitrate:         "192k",
					AudioSampleRate:      48000,
					Bitrate:              "5000k",
					Framerate:            "30",
					KeepAspectRatio:      YesNoBoolean(false),
					SetAspectRatio:       "1280 x 1080",
					VideoCodec:           "libvpx",
					VideoCodecParameters: "no",
					Size:                 "1920x1080",
				},
			},
			{
				Name:   "webm_480p",
				Type:   "user",
				Output: "webm",
				Format: PresetFormat{
					Output:               "webm",
					AudioBitrate:         "192k",
					AudioSampleRate:      48000,
					Bitrate:              "2500k",
					Framerate:            "30",
					VideoCodec:           "libvpx",
					VideoCodecParameters: "no",
					Size:                 "1280x480",
				},
			},
		},
		UIPresets: []Preset{
			{
				Name:   "3GP 144p 256k",
				Type:   "ui",
				Output: "3gp",
				Format: PresetFormat{
					Output:              "3gp",
					Size:                "176x144",
					Bitrate:             "256k",
					AudioBitrate:        "12.2k",
					AudioSampleRate:     8000,
					AudioChannelsNumber: "1",
					Framerate:           "15",
					KeepAspectRatio:     YesNoBoolean(true),
					VideoCodec:          "h263",
					Profile:             "baseline",
					AudioCodec:          "libamr_nb",
					TwoPass:             YesNoBoolean(false),
					Turbo:               YesNoBoolean(false),
					TwinTurbo:           YesNoBoolean(false),
					Keyframe:            "90",
					AudioVolume:         100,
					Rotate:              "def",
					StripChapters:       YesNoBoolean(false),
				},
			},
			{
				Name:   "3GP 288p 360k",
				Type:   "ui",
				Output: "3gp",
				Format: PresetFormat{
					Output:              "3gp",
					Size:                "352x288",
					Bitrate:             "360k",
					AudioBitrate:        "12.2k",
					AudioSampleRate:     8000,
					AudioChannelsNumber: "1",
					Framerate:           "24",
					KeepAspectRatio:     YesNoBoolean(true),
					VideoCodec:          "h263",
					Profile:             "baseline",
					AudioCodec:          "libamr_nb",
					TwoPass:             YesNoBoolean(false),
					Turbo:               YesNoBoolean(false),
					TwinTurbo:           YesNoBoolean(false),
					Keyframe:            "120",
					AudioVolume:         100,
					Rotate:              "def",
					StripChapters:       YesNoBoolean(false),
				},
			},
		},
	})
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "GetPresetsList")
	c.Assert(req.query["type"], check.Equals, string(AllPresets))
}

func (s *S) TestGetPreset(c *check.C) {
	server, requests := s.startServer(`
{
	"response": {
		"name":"webm_1080p",
		"type":"user",
		"output":"webm",
		"format":{
			"output":"webm",
			"audio_bitrate":"192k",
			"audio_sample_rate":"48000",
			"bitrate":"5000k",
			"framerate":"30",
			"keep_aspect_ratio":"no",
			"set_aspect_ratio":"1280 x 1080",
			"video_codec":"libvpx",
			"video_codec_parameters":"no",
			"size":"1920x1080"
		}
	}
}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	preset, err := client.GetPreset("webm_1080p")
	c.Assert(err, check.IsNil)
	expected := Preset{
		Name:   "webm_1080p",
		Type:   "user",
		Output: "webm",
		Format: PresetFormat{
			Output:               "webm",
			AudioBitrate:         "192k",
			AudioSampleRate:      48000,
			Bitrate:              "5000k",
			Framerate:            "30",
			KeepAspectRatio:      YesNoBoolean(false),
			SetAspectRatio:       "1280 x 1080",
			VideoCodec:           "libvpx",
			VideoCodecParameters: "no",
			Size:                 "1920x1080",
		},
	}
	c.Assert(*preset, check.DeepEquals, expected)

	req := <-requests
	c.Assert(req.query["action"], check.Equals, "GetPreset")
	c.Assert(req.query["type"], check.Equals, string(AllPresets))
	c.Assert(req.query["name"], check.Equals, "webm_1080p")
}

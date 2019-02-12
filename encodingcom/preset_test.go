package encodingcom

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestGetPresetsList(t *testing.T) {
	server, requests := startServer(`
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
		},
		{
			"name":"sample_hls",
			"type":"user",
			"output":"advanced_hls",
			"format":{
				"output":"advanced_hls",
				"stream":{
					"audio_bitrate":"64k",
					"audio_codec":"dolby_aac",
					"audio_volume":"100",
					"bitrate":"1000k",
					"keyframe":"90",
					"profile":"Main",
					"size":"1080x720",
					"two_pass":"yes",
					"video_codec":"libx264",
					"video_codec_parameters": {
						"keyint_min": "25",
						"sc_threshold": "40"
					}
				}
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
	if err != nil {
		t.Fatal(err)
	}
	expectedResp := &ListPresetsResponse{
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
			{
				Name:   "sample_hls",
				Type:   "user",
				Output: "advanced_hls",
				Format: PresetFormat{
					Output: "advanced_hls",
					StreamRawMap: map[string]interface{}{
						"audio_bitrate": "64k",
						"audio_volume":  "100",
						"size":          "1080x720",
						"two_pass":      "yes",
						"video_codec":   "libx264",
						"audio_codec":   "dolby_aac",
						"bitrate":       "1000k",
						"keyframe":      "90",
						"profile":       "Main",
						"video_codec_parameters": map[string]interface{}{
							"keyint_min":   "25",
							"sc_threshold": "40",
						},
					},
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
	}

	if !reflect.DeepEqual(resp, expectedResp) {
		t.Errorf("wrong response returned\nwant %#v\ngot %#v", expectedResp, resp)
	}

	expectedQuery := map[string]interface{}{
		"action":  "GetPresetsList",
		"type":    string(AllPresets),
		"userid":  "myuser",
		"userkey": "123",
	}
	req := <-requests
	if !reflect.DeepEqual(req.query, expectedQuery) {
		t.Errorf("wrong query sent to encoding.com\nwant %#v\ngot  %#v", expectedQuery, req.query)
	}

	sampleHlsStream := Stream{
		AudioBitrate: "64k",
		AudioVolume:  100,
		Size:         "1080x720",
		TwoPass:      YesNoBoolean(true),
		VideoCodec:   "libx264",
		AudioCodec:   "dolby_aac",
		Bitrate:      "1000k",
		Keyframe:     "90",
		Profile:      "Main",
		VideoCodecParametersRaw: map[string]interface{}{
			"sc_threshold": "40",
			"keyint_min":   "25",
		},
	}
	sampleVideoCodecParams := VideoCodecParameters{
		ScThreshold: "40",
		KeyIntMin:   "25",
	}
	for _, hlsPreset := range resp.UserPresets {
		if hlsPreset.Output == "advanced_hls" {
			streams := hlsPreset.Format.Stream()
			if !reflect.DeepEqual(streams[0], sampleHlsStream) {
				t.Errorf("wrong hls stream\nwant %#v\ngot  %#v", sampleHlsStream, streams[0])
			}

			if params := streams[0].VideoCodecParameters(); !reflect.DeepEqual(params, sampleVideoCodecParams) {
				t.Errorf("wrong video codec params for hls stream\nwant %#v\ngot  %#v", sampleVideoCodecParams, params)
			}
		}
	}
}

func TestListPresetsError(t *testing.T) {
	server, _ := startServer(`{"response": {"errors": {"error": "who moved my preset?"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.ListPresets(AllPresets)
	if err == nil {
		t.Error("unexpect <nil> error")
	}
	if resp != nil {
		t.Errorf("unexpected non-nil resp: %#v", resp)
	}
}

func TestPresetFormatStream(t *testing.T) {
	var tests = []struct {
		testCase  string
		streamRaw interface{}
		expected  []Stream
	}{
		{
			"single stream",
			map[string]interface{}{
				"audio_bitrate": "64k",
				"audio_volume":  "100",
				"size":          "1080x720",
				"two_pass":      "yes",
			},
			[]Stream{
				{
					AudioBitrate: "64k",
					AudioVolume:  100,
					Size:         "1080x720",
					TwoPass:      YesNoBoolean(true),
				},
			},
		},
		{
			"multiple streams",
			[]map[string]interface{}{
				{
					"audio_bitrate": "64k",
					"audio_volume":  "100",
					"size":          "1080x720",
					"two_pass":      "yes",
				},
				{
					"audio_bitrate": "128k",
					"audio_volume":  "100",
					"size":          "1920x1080",
					"two_pass":      "yes",
				},
			},
			[]Stream{
				{
					AudioBitrate: "64k",
					AudioVolume:  100,
					Size:         "1080x720",
					TwoPass:      YesNoBoolean(true),
				},
				{
					AudioBitrate: "128k",
					AudioVolume:  100,
					Size:         "1920x1080",
					TwoPass:      YesNoBoolean(true),
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, func(t *testing.T) {
			p := PresetFormat{StreamRawMap: test.streamRaw}
			streams := p.Stream()
			if !reflect.DeepEqual(streams, test.expected) {
				t.Errorf("wrong streams\nwant %#v\ngot  %#v", test.expected, streams)
			}
		})
	}
}

func TestGetPreset(t *testing.T) {
	server, requests := startServer(`
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
	if err != nil {
		t.Fatal(err)
	}
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
	if !reflect.DeepEqual(*preset, expected) {
		t.Errorf("wrong preset returned\nwant %#v\ngot  %#v", expected, *preset)
	}

	expectedQuery := map[string]interface{}{
		"action":  "GetPreset",
		"type":    string(AllPresets),
		"name":    "webm_1080p",
		"userid":  "myuser",
		"userkey": "123",
	}
	req := <-requests
	if !reflect.DeepEqual(req.query, expectedQuery) {
		t.Errorf("wrong query sent to encoding.com\nwant %#v\ngot  %#v", expectedQuery, req.query)
	}
}

func TestGetPresetError(t *testing.T) {
	server, _ := startServer(`{"response": {"errors": {"error": "can't get no presetisfaction"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.GetPreset("my-preciouset")
	if err == nil {
		t.Error("unexpect <nil> error")
	}
	if resp != nil {
		t.Errorf("unexpected non-nil resp: %#v", resp)
	}
}

func TestSavePreset(t *testing.T) {
	server, requests := startServer(`
	{
		"response": {
			"message": "Saved",
			"SavedPreset": "mp4_1080p"
		}
	}
`)
	defer server.Close()

	const presetName = "mp4_1080p"
	format := Format{
		VideoCodec:   "x264",
		AudioCodec:   "aac",
		Bitrate:      "900k",
		AudioBitrate: "64k",
		Size:         "1920x1080",
	}
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.SavePreset(presetName, format)
	if err != nil {
		t.Fatal(err)
	}
	if resp.SavedPreset != presetName {
		t.Errorf("wrong preset name returned\nwant %q\ngot  %q", presetName, resp.SavedPreset)
	}

	rawFormat, err := json.Marshal([]Format{format})
	if err != nil {
		t.Fatal(err)
	}
	var expectedFormat []interface{}
	err = json.Unmarshal(rawFormat, &expectedFormat)
	if err != nil {
		t.Fatal(err)
	}

	expectedQuery := map[string]interface{}{
		"action":  "SavePreset",
		"name":    presetName,
		"format":  expectedFormat,
		"userid":  "myuser",
		"userkey": "123",
	}
	req := <-requests
	if !reflect.DeepEqual(req.query, expectedQuery) {
		t.Errorf("wrong query sent to encoding.com\nwant %#v\ngot  %#v", expectedQuery, req.query)
	}
}

func TestSavePresetError(t *testing.T) {
	server, _ := startServer(`{"response": {"errors": {"error": "incomplete preset data"}}}`)
	defer server.Close()

	format := Format{VideoCodec: "x264"}
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.SavePreset("preset-1", format)
	if err == nil {
		t.Error("unexpect <nil> error")
	}
	if resp != nil {
		t.Errorf("unexpected non-nil resp: %#v", resp)
	}
}

func TestDeletePreset(t *testing.T) {
	server, requests := startServer(`
	{
		"response":{
			"message":"Deleted"
		}
	}
`)
	defer server.Close()

	const presetName = "mp4_1080p"
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.DeletePreset(presetName)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Message != "Deleted" {
		t.Errorf("wrong response message\nwant %q\ngot  %q", "Deleted", resp.Message)
	}

	expectedQuery := map[string]interface{}{
		"action":  "DeletePreset",
		"name":    presetName,
		"userid":  "myuser",
		"userkey": "123",
	}
	req := <-requests
	if !reflect.DeepEqual(req.query, expectedQuery) {
		t.Errorf("wrong query sent to encoding.com\nwant %#v\ngot  %#v", expectedQuery, req.query)
	}
}

func TestDeletePresetError(t *testing.T) {
	server, _ := startServer(`{"response": {"errors": {"error": "no preset, try postset"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.DeletePreset("some-preset")
	if err == nil {
		t.Error("unexpect <nil> error")
	}
	if resp != nil {
		t.Errorf("unexpected non-nil resp: %#v", resp)
	}
}

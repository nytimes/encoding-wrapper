package encodingdotcom

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	Endpoint string
}

func (c *Client) do(r *Request) (*http.Response, error) {
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	rawMsg := json.RawMessage(jsonRequest)
	m := map[string]interface{}{"query": &rawMsg}
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("json", string(data))
	req, err := http.NewRequest("POST", c.Endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return http.DefaultClient.Do(req)
}

type Request struct {
	UserID                  string       `json:"userid"`
	UserKey                 string       `json:"userkey"`
	Action                  string       `json:"action"`
	MediaID                 string       `json:"mediaid"`
	Source                  []string     `json:"source"`
	SplitScreen             *SplitScreen `json:"split_screen,omitempty"`
	Region                  string       `json:"region,omitempty"`
	NotifyFormat            string       `json:"notify_format,omitempty"`
	NotifyURL               string       `json:"notify,omitempty"`
	NotifyEncodingErrorsURL string       `json:"notify_encoding_errors,omitempty"`
	NotifyUploadURL         string       `json:"notify_upload,omitempty"`
	Format                  *Format      `json:"format,omitempty"`
}

type SplitScreen struct {
	Columns       int `json:"columns,string,omitempty"`
	Rows          int `json:"rows,string,omitempty"`
	PaddingLeft   int `json:"padding_left,string,omitempty"`
	PaddingRight  int `json:"padding_right,string,omitempty"`
	PaddingBottom int `json:"padding_bottom,string,omitempty"`
	PaddingTop    int `json:"padding_top,string,omitempty"`
}

type AddMediaResponse struct {
	Message string `json:"message,omitempty"`
	MediaID string `json:"mediaid,omitempty"`
}

type GetMediaListResponse struct {
	Media []struct {
		MediaFile   string    `json:"mediafile,omitempty"`
		MediaID     string    `json:"mediaid,omitempty"`
		MediaStatus string    `json:"mediastatus,omitempty"`
		CreatedDate time.Time `json:"createdate,string,omitempty"`
		StartDate   time.Time `json:"startdate,string,omitempty"`
		FinishDate  time.Time `json:"finishdate,string,omitempty"`
	}
}

type Format struct {
	NoiseReduction          string        `json:"noise_reduction,omitempty"`
	Output                  []string      `json:"output,omitempty"`
	VideoCodec              string        `json:"video_codec,omitempty"`
	AudioCodec              string        `json:"audio_codec,omitempty"`
	Bitrate                 string        `json:"bitrate,omitempty"`
	AudioBitrate            string        `json:"audio_bitrate,omitempty"`
	AudioSampleRate         string        `json:"audio_sample_rate,omitempty"`
	AudioChannelsNumber     int           `json:"audio_channels_number,string,omitempty"`
	AudioVolume             uint          `json:"audio_volume,string,omitempty"`
	Framerate               string        `json:"framerate,omitempty"`
	FramerateUpperThreshold string        `json:"framerate_upper_threshold,omitempty"`
	Size                    string        `json:"size,omitempty"`
	FadeIn                  string        `json:"fade_in,omitempty"`
	FadeOut                 string        `json:"fade_out,omitempty"`
	CropLeft                int           `json:"crop_left,string,omitempty"`
	CropTop                 int           `json:"crop_top,string,omitempty"`
	CropRight               int           `json:"crop_right,string,omitempty"`
	CropBottom              int           `json:"crop_bottom,string,omitempty"`
	KeepAspectRatio         YesNoBoolean  `json:"keep_aspect_ratio,omitempty"`
	SetAspectRatio          string        `json:"set_aspect_ratio,omitempty"`
	AddMeta                 YesNoBoolean  `json:"add_meta,omitempty"`
	Hint                    YesNoBoolean  `json:"hint,omitempty"`
	RcInitOccupancy         string        `json:"rc_init_occupancy,omitempty"`
	MinRate                 string        `json:"minrate,omitempty"`
	MaxRate                 string        `json:"maxrate,omitempty"`
	BufSize                 string        `json:"bufsize,omitempty"`
	Keyframe                []string      `json:"keyframe,omitempty"`
	Start                   string        `json:"start,omitempty"`
	Duration                string        `json:"duration,omitempty"`
	ForceKeyframes          string        `json:"force_keyframes,omitempty"`
	Bframes                 int           `json:"bframes,string,omitempty"`
	Gop                     string        `json:"gop,omitempty"`
	Metadata                *Metadata     `json:"metadata,omitempty"`
	Destination             []string      `json:"destination,omitempty"`
	Logo                    *Logo         `json:"logo,omitempty"`
	Overlay                 []Overlay     `json:"overlay,omitempty"`
	TextOverlay             []TextOverlay `json:"text_overlay,omitempty"`
	VideoCodecParameters    string        `json:"video_codec_parameters,omitempty"`
	Profile                 string        `json:"profile,omitempty"`
	Turbo                   string        `json:"turbo,omitempty"`
	Rotate                  string        `json:"rotate,omitempty"`
	SetRotate               string        `json:"set_rotate,omitempty"`
	AudioSync               string        `json:"audio_sync,omitempty"`
	VideoSync               string        `json:"video_sync,omitempty"`
	ForceInterlaced         string        `json:"force_interlaced,omitempty"`
	StripChapters           YesNoBoolean  `json:"strip_chapters,omitempty"`
}

type Logo struct {
	LogoSourceURL string `json:"logo_source,omitempty"`
	LogoX         int    `json:"logo_x,string,omitempty"`
	LogoY         int    `json:"logo_y,string,omitempty"`
	LogoMode      int    `json:"logo_mode,string,omitempty"`
	LogoThreshold string `json:"logo_threshold,omitempty"`
}

type Overlay struct {
	OverlaySource   string  `json:"overlay_source,omitempty"`
	OverlayLeft     string  `json:"overlay_left,omitempty"`
	OverlayRight    string  `json:"overlay_right,omitempty"`
	OverlayTop      string  `json:"overlay_top,omitempty"`
	OverlayBottom   string  `json:"overlay_bottom`
	Size            string  `json:"size,omitempty"`
	OverlayStart    float64 `json:"overlay_start,string,omitempty"`
	OverlayDuration float64 `json:"overlay_duration,string,omitempty"`
}

type TextOverlay struct {
	Text            []string       `json:"text,omitempty"`
	FontSourceURL   string         `json:"font_source,omitempty"`
	FontSize        uint           `json:"font_size,string,omitempty"`
	FontRotate      int            `json:"font_rotate,string,omitempty"`
	FontColor       string         `json:"font_color,omitempty"`
	AlignCenter     ZeroOneBoolean `json:"align_center,omitempty"`
	OverlayX        int            `json:"overlay_x,string,omitempty"`
	OverlayY        int            `json:"overlay_y,string,omitempty"`
	Size            string         `json:"size,omitempty"`
	OverlayStart    float64        `json:"overlay_start,string,omitempty"`
	OverlayDuration float64        `json:"overlay_duration,string,omitempty"`
}

type Metadata struct {
	Title       string `json:"title,omitempty"`
	Copyright   string `json:"copyright,omitempty"`
	Author      string `json:"author,omitempty"`
	Description string `json:"description,omitempty"`
	Album       string `json:"album,omitempty"`
}

type YesNoBoolean bool
type ZeroOneBoolean bool

func (b YesNoBoolean) MarshalJSON() ([]byte, error) {
	return boolToByte(bool(b), "yes", "no"), nil
}

func (b ZeroOneBoolean) MarshalJSON() ([]byte, error) {
	return boolToByte(bool(b), "1", "0"), nil
}

func boolToByte(b bool, t, f string) []byte {
	if b {
		return []byte(`"` + t + `"`)
	}
	return []byte(`"` + f + `"`)
}

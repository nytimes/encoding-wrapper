// Package encodingcom provides types and methods for interacting with the
// Encoding.com API.
//
// You can get more details on the API at http://api.encoding.com/.
package encodingcom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Client is the basic type for interacting with the API. It provides methods
// matching the available actions in the API.
type Client struct {
	Endpoint string
	UserID   string
	UserKey  string
}

func (c *Client) do(r *request, out interface{}) error {
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return err
	}
	rawMsg := json.RawMessage(jsonRequest)
	m := map[string]interface{}{"query": &rawMsg}
	reqData, err := json.Marshal(m)
	if err != nil {
		return err
	}
	params := url.Values{}
	params.Add("json", string(reqData))
	req, err := http.NewRequest("POST", c.Endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var errRespWrapper map[string]*errorResponse
	err = json.Unmarshal(respData, &errRespWrapper)
	if err != nil {
		return err
	}
	if errResp := errRespWrapper["response"]; errResp.Errors != nil {
		return &APIError{
			Message: errResp.Message,
			Errors:  errResp.Errors.Error,
		}
	}
	return json.Unmarshal(respData, out)
}

// APIError represents an error returned by the Encoding.com API.
//
// See http://goo.gl/BzvXZt for more details.
type APIError struct {
	Message string `json:",omitempty"`
	Errors  []string
}

// Error converts the whole interlying information to a representative string.
//
// It encodes the list of errors in JSON format.
func (apiErr *APIError) Error() string {
	data, _ := json.Marshal(apiErr)
	return fmt.Sprintf("Error returned by the Encoding.com API: %s", data)
}

type errorResponse struct {
	Message string  `json:"message,omitempty"`
	Errors  *errors `json:"errors,omitempty"`
}

type errors struct {
	Error []string `json:"error,omitempty"`
}

type request struct {
	UserID                  string       `json:"userid"`
	UserKey                 string       `json:"userkey"`
	Action                  string       `json:"action"`
	MediaID                 string       `json:"mediaid,omitempty"`
	Source                  []string     `json:"source,omitempty"`
	SplitScreen             *SplitScreen `json:"split_screen,omitempty"`
	Region                  string       `json:"region,omitempty"`
	NotifyFormat            string       `json:"notify_format,omitempty"`
	NotifyURL               string       `json:"notify,omitempty"`
	NotifyEncodingErrorsURL string       `json:"notify_encoding_errors,omitempty"`
	NotifyUploadURL         string       `json:"notify_upload,omitempty"`
	Format                  *Format      `json:"format,omitempty"`
}

// SplitScreen is the set of options for combining several sources to one split
// screen video.
//
// See http://goo.gl/EolKyv for more details.
type SplitScreen struct {
	Columns       int `json:"columns,string,omitempty"`
	Rows          int `json:"rows,string,omitempty"`
	PaddingLeft   int `json:"padding_left,string,omitempty"`
	PaddingRight  int `json:"padding_right,string,omitempty"`
	PaddingBottom int `json:"padding_bottom,string,omitempty"`
	PaddingTop    int `json:"padding_top,string,omitempty"`
}

// Format is the set of options for defining the output format when encoding
// new media files.
//
// See http://goo.gl/dcE1pF for more details.
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

// Logo is the set of options for watermarking media during encoding, allowing
// users to add a image to the final media.
//
// See http://goo.gl/4z2Q5S for more details.
type Logo struct {
	LogoSourceURL string `json:"logo_source,omitempty"`
	LogoX         int    `json:"logo_x,string,omitempty"`
	LogoY         int    `json:"logo_y,string,omitempty"`
	LogoMode      int    `json:"logo_mode,string,omitempty"`
	LogoThreshold string `json:"logo_threshold,omitempty"`
}

// Overlay is the set of options for adding a video overlay in the media being
// encoded.
//
// See http://goo.gl/Q6sjkR for more details.
type Overlay struct {
	OverlaySource   string  `json:"overlay_source,omitempty"`
	OverlayLeft     string  `json:"overlay_left,omitempty"`
	OverlayRight    string  `json:"overlay_right,omitempty"`
	OverlayTop      string  `json:"overlay_top,omitempty"`
	OverlayBottom   string  `json:"overlay_bottom"`
	Size            string  `json:"size,omitempty"`
	OverlayStart    float64 `json:"overlay_start,string,omitempty"`
	OverlayDuration float64 `json:"overlay_duration,string,omitempty"`
}

// TextOverlay is the set of options for adding a text overlay in the media
// being encoded.
//
// See http://goo.gl/gUKi5t for more details.
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

// Metadata represents media metadata, as provided in the Format struct when
// encoding new media.
//
// See http://goo.gl/jNSio9 for more details.
type Metadata struct {
	Title       string `json:"title,omitempty"`
	Copyright   string `json:"copyright,omitempty"`
	Author      string `json:"author,omitempty"`
	Description string `json:"description,omitempty"`
	Album       string `json:"album,omitempty"`
}

// YesNoBoolean is a boolean that turns true into "yes" and false into "no"
// when encoded as JSON.
type YesNoBoolean bool

// MarshalJSON is the method that ensures that YesNoBoolean satisfies the
// json.Marshaler interface.
func (b YesNoBoolean) MarshalJSON() ([]byte, error) {
	return boolToBytes(bool(b), "yes", "no"), nil
}

// ZeroOneBoolean is a boolean that turns true into "1" and false into "0" when
// encoded as JSON.
type ZeroOneBoolean bool

// MarshalJSON is the method that ensures that ZeroOneBoolean satisfies the
// json.Marshaler interface.
func (b ZeroOneBoolean) MarshalJSON() ([]byte, error) {
	return boolToBytes(bool(b), "1", "0"), nil
}

func boolToBytes(b bool, t, f string) []byte {
	if b {
		return []byte(`"` + t + `"`)
	}
	return []byte(`"` + f + `"`)
}

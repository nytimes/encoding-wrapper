package encodingdotcom

import "time"

type Client struct {
	Endpoint string
}

type Action string

const (
	AddMedia Action = "AddMedia"
)

type Request struct {
	UserID                  string
	UserKey                 string
	Action                  Action
	MediaID                 string
	Source                  []string
	SplitScreen             SplitScreen
	Region                  string
	NotifyFormat            string
	NotifyURL               string `json:"notify"`
	NotifyEncodingErrorsURL string `json:"notify_encoding_errors"`
	NotifyUploadURL         string `json:"notify_upload"`
	Format                  Format
}

type SplitScreen struct {
	Columns       int `json:"columns",string`
	Rows          int `json:"rows",string`
	PaddingLeft   int `json:"padding_left",string`
	PaddingRight  int `json:"padding_right",string`
	PaddingBottom int `json:"padding_bottom",string`
	PaddingTop    int `json:"padding_top",string`
}

type AddMediaResponse struct {
	Message string
	MediaID string
}

type GetMediaListResponse struct {
	Media []struct {
		MediaFile   string
		MediaID     string
		MediaStatus string
		CreatedDate time.Time
		StartDate   time.Time
		FinishDate  time.Time
	}
}

type Format struct {
	NoiseReduction          string
	Output                  []string
	VideoCodec              string
	AudioCodec              string
	Bitrate                 string
	AudioBitrate            string
	AudioSampleRate         string
	AudioChannelsNumber     int  `json:"audio_channels_number",string`
	AudioVolume             uint `json:"audio_volume",string`
	Framerate               string
	FramerateUpperThreshold string
	Size                    string
	FadeIn                  string
	FadeOut                 string
	CropLeft                int `json:"crop_left",string`
	CropTop                 int `json:"crop_top",string`
	CropRight               int `json:"crop_right",string`
	CropBottom              int `json:"crop_bottom",string`
	KeepAspectRatio         YesNoBoolean
	SetAspectRatio          string
	AddMeta                 YesNoBoolean
	Hint                    YesNoBoolean
	RcInitOccupancy         string
	MinRate                 string
	MaxRate                 string
	BufSize                 string
	Keyframe                []string
	Start                   string
	Duration                string
	ForceKeyframes          string
	Bframes                 int `json:"bframes",string`
	Gop                     string
	Metadata                Metadata
	Destination             []string
	Logo                    Logo
	Overlay                 []Overlay
	TextOverlay             []TextOverlay
	VideoCodecParameters    string
	Profile                 string
	Turbo                   string
	Rotate                  string
	SetRotate               string
	AudioSync               string
	VideoSync               string
	ForceInterlaced         string
	StripChapters           YesNoBoolean
}

type Logo struct {
	LogoSourceURL string `json:"logo_source"`
	LogoX         int    `json:"logo_x",string`
	LogoY         int    `json:"logo_y",string`
	LogoMode      string
	LogoThreshold string
}

type Overlay struct {
	OverlaySource   string
	OverlayLeft     string
	OverlayRight    string
	OverlayTop      string
	OverlayBottom   string
	Size            string
	OverlayStart    float64 `json:"overlay_start",string`
	OverlayDuration float64 `json:"overlay_duration",string`
}

type TextOverlay struct {
	Text            []string
	FontSource      string
	FontSize        string
	FontRotate      string
	FontColor       string
	AlignCenter     ZeroOneBoolean
	OverlayX        int `json:"overlay_x",string`
	OverlayY        int `json:"overlay_y",string`
	Size            string
	OverlayStart    float64 `json:"overlay_start",string`
	OverlayDuration float64 `json:"overlay_duration",string`
}

type Metadata struct {
	Title       string
	Copyright   string
	Author      string
	Description string
	Album       string
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

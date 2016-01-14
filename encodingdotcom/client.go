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
	UserID                  string      `json:"userid"`
	UserKey                 string      `json:"userkey"`
	Action                  Action      `json:"action"`
	MediaID                 string      `json:"mediaid"`
	Source                  []string    `json:"source"`
	SplitScreen             SplitScreen `json:"split_screen"`
	Region                  string      `json:"region"`
	NotifyFormat            string      `json:"notify_format"`
	NotifyURL               string      `json:"notify"`
	NotifyEncodingErrorsURL string      `json:"notify_encoding_errors"`
	NotifyUploadURL         string      `json:"notify_upload"`
	Format                  Format      `json:"format"`
}

type SplitScreen struct {
	Columns       int `json:"columns,string"`
	Rows          int `json:"rows,string"`
	PaddingLeft   int `json:"padding_left,string"`
	PaddingRight  int `json:"padding_right,string"`
	PaddingBottom int `json:"padding_bottom,string"`
	PaddingTop    int `json:"padding_top,string"`
}

type AddMediaResponse struct {
	Message string `json:"message"`
	MediaID string `json:"mediaid"`
}

type GetMediaListResponse struct {
	Media []struct {
		MediaFile   string    `json:"mediafile"`
		MediaID     string    `json:"mediaid"`
		MediaStatus string    `json:"mediastatus"`
		CreatedDate time.Time `json:"createdate,string"`
		StartDate   time.Time `json:"startdate,string"`
		FinishDate  time.Time `json:"finishdate,string"`
	}
}

type Format struct {
	NoiseReduction          string        `json:"noise_reduction"`
	Output                  []string      `json:"output"`
	VideoCodec              string        `json:"video_codec"`
	AudioCodec              string        `json:"audio_codec"`
	Bitrate                 string        `json:"bitrate"`
	AudioBitrate            string        `json:"audio_bitrate"`
	AudioSampleRate         string        `json:"audio_sample_rate"`
	AudioChannelsNumber     int           `json:"audio_channels_number,string"`
	AudioVolume             uint          `json:"audio_volume,string"`
	Framerate               string        `json:"framerate"`
	FramerateUpperThreshold string        `json:"framerate_upper_threshold"`
	Size                    string        `json:"size"`
	FadeIn                  string        `json:"fade_in"`
	FadeOut                 string        `json:"fade_out"`
	CropLeft                int           `json:"crop_left,string"`
	CropTop                 int           `json:"crop_top,string"`
	CropRight               int           `json:"crop_right,string"`
	CropBottom              int           `json:"crop_bottom,string"`
	KeepAspectRatio         YesNoBoolean  `json:"keep_aspect_ratio"`
	SetAspectRatio          string        `json:"set_aspect_ratio"`
	AddMeta                 YesNoBoolean  `json:"add_meta"`
	Hint                    YesNoBoolean  `json:"hint"`
	RcInitOccupancy         string        `json:"rc_init_occupancy"`
	MinRate                 string        `json:"minrate"`
	MaxRate                 string        `json:"maxrate"`
	BufSize                 string        `json:"bufsize"`
	Keyframe                []string      `json:"keyframe"`
	Start                   string        `json:"start"`
	Duration                string        `json:"duration"`
	ForceKeyframes          string        `json:"force_keyframes"`
	Bframes                 int           `json:"bframes,string"`
	Gop                     string        `json:"gop"`
	Metadata                Metadata      `json:"metadata"`
	Destination             []string      `json:"destination"`
	Logo                    Logo          `json:"logo"`
	Overlay                 []Overlay     `json:"overlay"`
	TextOverlay             []TextOverlay `json:"text_overlay"`
	VideoCodecParameters    string        `json:"video_codec_parameters"`
	Profile                 string        `json:"profile"`
	Turbo                   string        `json:"turbo"`
	Rotate                  string        `json:"rotate"`
	SetRotate               string        `json:"set_rotate"`
	AudioSync               string        `json:"audio_sync"`
	VideoSync               string        `json:"video_sync"`
	ForceInterlaced         string        `json:"force_interlaced"`
	StripChapters           YesNoBoolean  `json:"strip_chapters"`
}

type Logo struct {
	LogoSourceURL string `json:"logo_source"`
	LogoX         int    `json:"logo_x,string"`
	LogoY         int    `json:"logo_y,string"`
	LogoMode      int    `json:"logo_mode,string"`
	LogoThreshold string `json:"logo_threshold"`
}

type Overlay struct {
	OverlaySource   string  `json:"overlay_source"`
	OverlayLeft     string  `json:"overlay_left"`
	OverlayRight    string  `json:"overlay_right"`
	OverlayTop      string  `json:"overlay_top"`
	OverlayBottom   string  `json:"overlay_bottom`
	Size            string  `json:"size"`
	OverlayStart    float64 `json:"overlay_start,string"`
	OverlayDuration float64 `json:"overlay_duration,string"`
}

type TextOverlay struct {
	Text            []string       `json:"text"`
	FontSourceURL   string         `json:"font_source"`
	FontSize        uint           `json:"font_size,string"`
	FontRotate      int            `json:"font_rotate,string"`
	FontColor       string         `json:"font_color"`
	AlignCenter     ZeroOneBoolean `json:"align_center"`
	OverlayX        int            `json:"overlay_x,string"`
	OverlayY        int            `json:"overlay_y,string"`
	Size            string         `json:"size"`
	OverlayStart    float64        `json:"overlay_start,string"`
	OverlayDuration float64        `json:"overlay_duration,string"`
}

type Metadata struct {
	Title       string `json:"title"`
	Copyright   string `json:"copyright"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Album       string `json:"album"`
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

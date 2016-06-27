package encodingcom

const (
	// AllPresets is used to retrieve all presets in the response of
	// ListPresets or GetPreset methods.
	AllPresets = PresetType("all")

	// UserPresets is used to retrieve only user-created presets in the
	// response of ListPresets or GetPreset methods.
	UserPresets = PresetType("user")

	// UIPresets is used to retrieve only ui (standard) presets in the
	// response of ListPresets or GetPreset methods.
	UIPresets = PresetType("ui")
)

// PresetType represents the type of preset used as the input of the
// ListPresets and GetPreset methods.
type PresetType string

// Preset represents a preset in the Encoding.com API.
type Preset struct {
	Name   string       `json:"name"`
	Type   PresetType   `json:"type"`
	Output string       `json:"output"`
	Format PresetFormat `json:"format"`
}

// PresetFormat is the set of options for defining the output format in
// presets.
type PresetFormat struct {
	NoiseReduction          string       `json:"noise_reduction,omitempty"`
	Output                  string       `json:"output,omitempty"`
	VideoCodec              string       `json:"video_codec,omitempty"`
	AudioCodec              string       `json:"audio_codec,omitempty"`
	Bitrate                 string       `json:"bitrate,omitempty"`
	AudioBitrate            string       `json:"audio_bitrate,omitempty"`
	AudioSampleRate         uint         `json:"audio_sample_rate,string,omitempty"`
	AudioChannelsNumber     string       `json:"audio_channels_number,omitempty"`
	AudioVolume             uint         `json:"audio_volume,string,omitempty"`
	Framerate               string       `json:"framerate,omitempty"`
	FramerateUpperThreshold string       `json:"framerate_upper_threshold,omitempty"`
	Size                    string       `json:"size,omitempty"`
	FadeIn                  string       `json:"fade_in,omitempty"`
	FadeOut                 string       `json:"fade_out,omitempty"`
	CropLeft                int          `json:"crop_left,string,omitempty"`
	CropTop                 int          `json:"crop_top,string,omitempty"`
	CropRight               int          `json:"crop_right,string,omitempty"`
	CropBottom              int          `json:"crop_bottom,string,omitempty"`
	KeepAspectRatio         YesNoBoolean `json:"keep_aspect_ratio,omitempty"`
	SetAspectRatio          string       `json:"set_aspect_ratio,omitempty"`
	AddMeta                 YesNoBoolean `json:"add_meta,omitempty"`
	Hint                    YesNoBoolean `json:"hint,omitempty"`
	RcInitOccupancy         string       `json:"rc_init_occupancy,omitempty"`
	MinRate                 string       `json:"minrate,omitempty"`
	MaxRate                 string       `json:"maxrate,omitempty"`
	BufSize                 string       `json:"bufsize,omitempty"`
	Keyframe                string       `json:"keyframe,omitempty"`
	Start                   string       `json:"start,omitempty"`
	Duration                string       `json:"duration,omitempty"`
	ForceKeyframes          string       `json:"force_keyframes,omitempty"`
	Bframes                 int          `json:"bframes,string,omitempty"`
	Gop                     string       `json:"gop,omitempty"`
	Metadata                *Metadata    `json:"metadata,omitempty"`
	SegmentDuration         uint         `json:"segment_duration,omitempty"`
	Logo                    *Logo        `json:"logo,omitempty"`
	VideoCodecParameters    string       `json:"video_codec_parameters,omitempty"`
	Profile                 string       `json:"profile,omitempty"`
	TwoPass                 YesNoBoolean `json:"two_pass,omitempty"`
	Turbo                   YesNoBoolean `json:"turbo,omitempty"`
	TwinTurbo               YesNoBoolean `json:"twin_turbo,omitempty"`
	Rotate                  string       `json:"rotate,omitempty"`
	SetRotate               string       `json:"set_rotate,omitempty"`
	AudioSync               string       `json:"audio_sync,omitempty"`
	VideoSync               string       `json:"video_sync,omitempty"`
	ForceInterlaced         string       `json:"force_interlaced,omitempty"`
	StripChapters           YesNoBoolean `json:"strip_chapters,omitempty"`
}

// ListPresetsResponse represents the response returned by the GetPresetsList
// action.
//
// See http://goo.gl/sugm5F for more details.
type ListPresetsResponse struct {
	UserPresets []Preset `json:"user"`
	UIPresets   []Preset `json:"ui"`
}

// ListPresets (GetPresetsList action in the Encoding.com API) returns a list
// of the presets matching the given type.
//
// See http://goo.gl/sugm5F for more details.
func (c *Client) ListPresets(presetType PresetType) (*ListPresetsResponse, error) {
	var result map[string]*ListPresetsResponse
	err := c.do(&request{Action: "GetPresetsList", Type: string(presetType)}, &result)
	if err != nil {
		return nil, err
	}
	return result["response"], nil
}

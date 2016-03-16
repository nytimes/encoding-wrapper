package elementalconductor

import (
	"net/http"

	"gopkg.in/check.v1"
)

func (s *S) TestGetPresets(c *check.C) {
	presetsResponseXML := `<?xml version="1.0" encoding="UTF-8"?>
<preset_list>
  <preset href="/presets/1" product="Elemental Conductor File + Audio Normalization Package + Audio Package" version="2.7.2vd.32545">
    <name>iPhone</name>
    <permalink>iphone</permalink>
    <description>Default output for iPhone</description>
    <preset_category href="/preset_categories/6">Devices</preset_category>
  </preset>
  <preset href="/presets/2" product="Elemental Conductor File + Audio Normalization Package + Audio Package" version="2.7.2vd.32545">
    <name>iPhone_ADAPT_HIGH</name>
    <permalink>iphone_adapt_high</permalink>
    <description>Default output for iPhone Adaptive high quality</description>
    <preset_category href="/preset_categories/6">Devices</preset_category>
  </preset>
  <next href="https://3e9n5rjaf3eb2.cloud.elementaltechnologies.com/presets?page=2&amp;amp;per_page=30"/>
</preset_list>`

	expectedPreset1 := Preset{
		Name:        "iPhone",
		Href:        "/presets/1",
		Permalink:   "iphone",
		Description: "Default output for iPhone",
	}

	expectedPreset2 := Preset{
		Name:        "iPhone_ADAPT_HIGH",
		Href:        "/presets/2",
		Permalink:   "iphone_adapt_high",
		Description: "Default output for iPhone Adaptive high quality",
	}

	var expectedOutput PresetList
	expectedOutput.Presets = make([]Preset, 2)
	expectedOutput.Presets[0] = expectedPreset1
	expectedOutput.Presets[1] = expectedPreset2

	server, _ := s.startServer(http.StatusCreated, presetsResponseXML)
	defer server.Close()

	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	getPresetsResponse, _ := client.GetPresets()
	c.Assert(getPresetsResponse, check.DeepEquals, &expectedOutput)
}

func (s *S) TestGetPreset(c *check.C) {
	presetResponseXML := `<?xml version="1.0" encoding="UTF-8"?>
<preset href="/presets/1" product="Elemental Conductor File + Audio Normalization Package + Audio Package" version="2.7.2vd.32545">
  <name>iPhone</name>
  <permalink>iphone</permalink>
  <description>Default output for iPhone</description>
  <preset_category href="/preset_categories/6">Devices</preset_category>
  <container>mp4</container>
  <mp4_settings>
    <id>1</id>
    <include_cslg>false</include_cslg>
    <mp4_major_brand nil="true"/>
    <progressive_downloading>true</progressive_downloading>
  </mp4_settings>
  <log_edit_points>false</log_edit_points>
  <video_description>
    <afd_signaling>None</afd_signaling>
    <anti_alias>true</anti_alias>
    <drop_frame_timecode>true</drop_frame_timecode>
    <encoder_type nil="true"/>
    <fixed_afd nil="true"/>
    <force_cpu_encode>false</force_cpu_encode>
    <height>320</height>
    <id>1</id>
    <insert_color_metadata>false</insert_color_metadata>
    <respond_to_afd>None</respond_to_afd>
    <sharpness>50</sharpness>
    <stretch_to_output>false</stretch_to_output>
    <timecode_passthrough>false</timecode_passthrough>
    <vbi_passthrough>false</vbi_passthrough>
    <width>480</width>
    <h264_settings>
      <adaptive_quantization>medium</adaptive_quantization>
      <bitrate>960000</bitrate>
      <buf_fill_pct nil="true"/>
      <buf_size nil="true"/>
      <cabac>false</cabac>
      <flicker_reduction>off</flicker_reduction>
      <force_field_pictures>false</force_field_pictures>
      <framerate_denominator>1</framerate_denominator>
      <framerate_follow_source>false</framerate_follow_source>
      <framerate_numerator>24</framerate_numerator>
      <gop_b_reference>false</gop_b_reference>
      <gop_closed_cadence>1</gop_closed_cadence>
      <gop_markers>false</gop_markers>
      <gop_num_b_frames>0</gop_num_b_frames>
      <gop_size>80</gop_size>
      <id>1</id>
      <interpolate_frc>false</interpolate_frc>
      <look_ahead_rate_control>medium</look_ahead_rate_control>
      <max_bitrate nil="true"/>
      <max_qp nil="true"/>
      <min_i_interval>0</min_i_interval>
      <min_qp nil="true"/>
      <num_ref_frames>1</num_ref_frames>
      <par_denominator>1</par_denominator>
      <par_follow_source>false</par_follow_source>
      <par_numerator>1</par_numerator>
      <passes>1</passes>
      <qp nil="true"/>
      <qp_step nil="true"/>
      <repeat_pps>false</repeat_pps>
      <scd>true</scd>
      <sei_timecode>false</sei_timecode>
      <slices>1</slices>
      <slow_pal>false</slow_pal>
      <softness nil="true"/>
      <svq>0</svq>
      <telecine>None</telecine>
      <transition_detection>false</transition_detection>
      <level>3</level>
      <profile>Baseline</profile>
      <rate_control_mode>ABR</rate_control_mode>
      <gop_mode>fixed</gop_mode>
      <interlace_mode>progressive</interlace_mode>
    </h264_settings>
    <gpu/>
    <selected_gpu nil="true"/>
    <codec>h.264</codec>
    <video_preprocessors>
      <deinterlacer>
        <algorithm>interpolate</algorithm>
        <deinterlace_mode>Deinterlace</deinterlace_mode>
        <force>false</force>
        <id>85</id>
      </deinterlacer>
    </video_preprocessors>
  </video_description>
  <audio_description>
    <audio_type>0</audio_type>
    <follow_input_audio_type>false</follow_input_audio_type>
    <follow_input_language_code>false</follow_input_language_code>
    <id>1</id>
    <language_code nil="true"/>
    <order>1</order>
    <stream_name nil="true"/>
    <aac_settings>
      <bitrate>128000</bitrate>
      <coding_mode>2_0</coding_mode>
      <id>1</id>
      <latm_loas>false</latm_loas>
      <mpeg2>false</mpeg2>
      <sample_rate>44100</sample_rate>
      <profile>LC</profile>
      <rate_control_mode>CBR</rate_control_mode>
    </aac_settings>
    <codec>aac</codec>
  </audio_description>
</preset>`

	expectedPreset := Preset{
		Name:          "iPhone",
		Href:          "/presets/1",
		Permalink:     "iphone",
		Description:   "Default output for iPhone",
		Container:     "mp4",
		VideoCodec:    "h.264",
		AudioCodec:    "aac",
		Width:         "480",
		Height:        "320",
		VideoBitrate:  "960000",
		AudioBitrate:  "128000",
		GopSize:       "80",
		GopMode:       "fixed",
		Profile:       "Baseline",
		ProfileLevel:  "3",
		RateControl:   "ABR",
		InterlaceMode: "progressive",
	}

	server, _ := s.startServer(http.StatusCreated, presetResponseXML)
	defer server.Close()

	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	getPresetResponse, _ := client.GetPreset("1")
	c.Assert(getPresetResponse, check.DeepEquals, &expectedPreset)
}

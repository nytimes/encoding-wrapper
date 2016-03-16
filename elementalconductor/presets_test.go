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
	c.Assert(getPresetsResponse, check.NotNil)
	c.Assert(getPresetsResponse, check.DeepEquals, &expectedOutput)
}

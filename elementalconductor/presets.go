package elementalconductor

// GetPresets returns a list of presets
func (c *Client) GetPresets() (*PresetList, error) {
	var result *PresetList
	err := c.do("GET", "/presets", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// PresetList represents the response returned by
// a query for the list of jobs
type PresetList struct {
	Presets []Preset `xml:"preset"`
}

// Preset represents a preset stored on Elemental Cloud
type Preset struct {
	Name        string `xml:"name"`
	Href        string `xml:"href,attr,omitempty"`
	Permalink   string `xml:"permalink,omitempty"`
	Description string `xml:"description,omitempty"`
}

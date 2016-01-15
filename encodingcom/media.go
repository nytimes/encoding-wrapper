package encodingcom

// AddMediaResponse represents the response returned by the AddMedia action.
//
// See http://goo.gl/Aqg8lc for more details.
type AddMediaResponse struct {
	Message string `json:"message,omitempty"`
	MediaID string `json:"mediaid,omitempty"`
}

// AddMedia adds a new media to user's queue.
//
// Format specifies details on how the source files are going to be encoded.
//
// See http://goo.gl/whvHwJ for more details on the source file formatting.
func (c *Client) AddMedia(source []string, format *Format) (*AddMediaResponse, error) {
	var result map[string]*AddMediaResponse
	err := c.do(&request{
		Action:  "AddMedia",
		Format:  format,
		Source:  source,
		UserID:  c.UserID,
		UserKey: c.UserKey,
	}, &result)
	if err != nil {
		return nil, err
	}
	return result["response"], nil
}

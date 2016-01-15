package encodingdotcom

func (c *Client) AddMedia(source []string, format *Format) (*AddMediaResponse, error) {
	var result map[string]*AddMediaResponse
	err := c.do(&Request{
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

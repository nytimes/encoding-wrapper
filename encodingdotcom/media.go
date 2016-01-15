package encodingdotcom

import (
	"encoding/json"
)

func (c *Client) AddMedia(source []string, format *Format) (*AddMediaResponse, error) {
	resp, err := c.do(&Request{
		Action:  "AddMedia",
		Format:  format,
		Source:  source,
		UserID:  c.UserID,
		UserKey: c.UserKey,
	})
	defer resp.Body.Close()
	var result map[string]*AddMediaResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result["response"], nil
}

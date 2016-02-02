// Package elementalcloud provides types and methods for interacting with the
// ElementalCloud.com API.
//
// You can get more details on the API at http://api.encoding.com/.
package elementalcloud

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client is the basic type for interacting with the API. It provides methods
// matching the available actions in the API.
type Client struct {
	Host    string
	UserID  string
	UserKey string
}

// NewClient creates a instance of the client type.
func NewClient(host, userID, userKey string) (*Client, error) {
	return &Client{Host: host, UserID: userID, UserKey: userKey}, nil
}

func (c *Client) do(method string, path string, body interface{}, out interface{}) error {
	xmlRequest, err := xml.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, c.Host+path, strings.NewReader(string(xmlRequest)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/xml")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return xml.Unmarshal(respData, out)
}

// FileInput contains location of the video file to be encoded
type FileInput struct {
	URI string `xml:"uri"`
}

// Job specifies the parameters for the Elemental Cloud job,
// where Profile is the id of an existing profile
type Job struct {
	FileInput FileInput `xml:"file_input"`
	Profile   string    `xml:"profile"`
}

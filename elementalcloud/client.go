// Package elementalcloud provides types and methods for interacting with the
// ElementalCloud.com API.
//
// You can get more details on the API at http://api.encoding.com/.
package elementalcloud

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Client is the basic type for interacting with the API. It provides methods
// matching the available actions in the API.
type Client struct {
	Host           string
	UserID         string
	APIKey         string
	ExpirationTime int
}

// NewClient creates a instance of the client type.
func NewClient(host, userID, apiKey string, expirationTime int) *Client {
	return &Client{Host: host, UserID: userID, APIKey: apiKey, ExpirationTime: expirationTime}
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
	req.Header.Set("Authorization", c.createAuthKey(path, c.ExpirationTime))
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

func (c *Client) createAuthKey(URL string, expireInSeconds int) string {
	now := time.Now()
	expire := string(now.Add(time.Duration(expireInSeconds) * time.Second).Unix())

	hasher := md5.New()
	hasher.Write([]byte(URL))
	hasher.Write([]byte(c.UserID))
	hasher.Write([]byte(c.APIKey))
	hasher.Write([]byte(expire))

	// innerKey := string((md5.Sum(URL + c.UserID + c.APIKey + expire)))
	return hex.EncodeToString(hasher.Sum(nil))
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

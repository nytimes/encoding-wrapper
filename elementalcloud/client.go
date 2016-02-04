// Package elementalcloud provides types and methods for interacting with the
// ElementalCloud API.
//
// You can get more details on the API at https://<elemental_server>/help/rest_api.
package elementalcloud

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Client is the basic type for interacting with the API. It provides methods
// matching the available actions in the API.
type Client struct {
	Host        string
	UserLogin   string
	APIKey      string
	AuthExpires int
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

// APIError represents an error returned by the Elemental Cloud REST API.
//
// See https://<elemental_server>/help/rest_api#rest_basics_errors_and_warnings
// for more details.
type APIError struct {
	Status int    `json:"status,omitempty"`
	Errors string `json:"errors,omitempty"`
}

// Error converts the whole interlying information to a representative string.
//
// It encodes the list of errors in JSON format.
func (apiErr *APIError) Error() string {
	data, _ := json.Marshal(apiErr)
	return fmt.Sprintf("Error returned by the Elemental Cloud REST Interface: %s", data)
}

// NewClient creates a instance of the client type.
func NewClient(host, userLogin, apiKey string, authExpires int) *Client {
	return &Client{Host: host, UserLogin: userLogin, APIKey: apiKey, AuthExpires: authExpires}
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
	expireTime := time.Now().Add(time.Duration(c.AuthExpires) * time.Second)
	req.Header.Set("Authorization", c.createAuthKey(path, expireTime))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return &APIError{
			Status: resp.StatusCode,
			Errors: string(respData),
		}
	}
	return xml.Unmarshal(respData, out)
}

func (c *Client) createAuthKey(URL string, expire time.Time) string {
	expireString := string(expire.Unix())
	hasher := md5.New()
	hasher.Write([]byte(URL))
	hasher.Write([]byte(c.UserLogin))
	hasher.Write([]byte(c.APIKey))
	hasher.Write([]byte(expireString))
	innerKey := hex.EncodeToString(hasher.Sum(nil))
	hasher = md5.New()
	hasher.Write([]byte(c.APIKey))
	hasher.Write([]byte(innerKey))
	return hex.EncodeToString(hasher.Sum(nil))
}

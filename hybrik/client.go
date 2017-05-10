package hybrik

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// APIInterface : Core interface for handling requests to the API
type APIInterface interface {
	connect() error
	isExpired() bool
	CallAPI(method string, apiPath string, params url.Values, body io.Reader) (string, error)
}

// ClientInterface : Wrappers around commonly used API calls
type ClientInterface interface {
	// Jobs API
	QueueJob(string) (string, error)
	GetJobInfo(string) (JobInfo, error)
	StopJob(string) error

	// Presets API
	GetPreset(string) (Preset, error)
	CreatePreset(Preset) (Preset, error)
	DeletePreset(string) error

	// Health Check
	HealthCheck() error
}

type Client struct {
	client APIInterface
}

// Config represents the configuration params that are necessary for the API calls to Hybrik to work
type Config struct {
	URL            string
	ComplianceDate string
	OAPIKey        string
	OAPISecret     string
	AuthKey        string
	AuthSecret     string
	OAPIURL        string
}

// API is the implementation of the HybrikAPI methods
type API struct {
	Config     Config
	token      string
	expiration string
}

type connectResponse struct {
	Token          string `json:"token"`
	ExpirationTime string `json:"expiration_time"`
}

// NewClient creates an instance of the HybrikAPI client
func NewClient(config Config) (*Client, error) {

	switch {
	case config.URL == "":
		return &Client{}, ErrNoAPIURL
	case config.OAPIKey == "":
		return &Client{}, ErrNoOAPIKey
	case config.OAPISecret == "":
		return &Client{}, ErrNoOAPISecret
	case config.AuthKey == "":
		return &Client{}, ErrNoAuthKey
	case config.AuthSecret == "":
		return &Client{}, ErrNoAuthSecret
	case config.ComplianceDate == "":
		return &Client{}, ErrNoComplianceDate
	case !regexp.MustCompile(`^\d{8}$`).MatchString(config.ComplianceDate):
		return &Client{}, ErrNoComplianceDate
	}

	_, err := url.ParseRequestURI(config.URL)
	if err != nil {
		return &Client{}, ErrInvalidURL
	}

	parts := strings.Split(config.URL, "//")
	if len(parts) < 2 {
		return &Client{}, ErrInvalidURL
	}

	config.OAPIURL = fmt.Sprintf("%s//%s:%s@%s",
		parts[0], config.OAPIKey, config.OAPISecret, parts[1],
	)

	return &Client{
		client: &API{
			Config: config,
		},
	}, nil
}

func (a *API) connect() error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", a.Config.OAPIURL), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Hybrik-Compliance", a.Config.ComplianceDate)

	query := req.URL.Query()
	query.Add("auth_key", a.Config.AuthKey)
	query.Add("auth_secret", a.Config.AuthSecret)
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var cr connectResponse
	err = json.NewDecoder(resp.Body).Decode(&cr)
	if err != nil {
		return err
	}

	a.token = cr.Token
	a.expiration = cr.ExpirationTime

	return nil
}

func (a API) isExpired() bool {
	if a.expiration == "" {
		return true
	}

	t, err := time.Parse(`2006-01-02T15:04:05.999Z`, a.expiration)
	if err != nil {
		return true
	}

	return time.Now().After(t)
}

// CallAPI is the general method to call for access to the API
func (a API) CallAPI(method string, apiPath string, params url.Values, body io.Reader) (string, error) {
	// Retrieves the 'token' and 'expiration_time'
	if a.isExpired() {
		err := a.connect()
		if err != nil {
			return "", err
		}
	}

	// Does the necessary http call here
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", a.Config.OAPIURL, apiPath), body)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Hybrik-Compliance", a.Config.ComplianceDate)
	req.Header.Set("X-Hybrik-Sapiauth", a.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.URL.RawQuery = params.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	callResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%d - %s %s: %s", resp.StatusCode, method, apiPath, string(callResp))
	}

	return string(callResp), nil
}

package elementalconductor

import (
	"encoding/xml"
	"strings"
	"time"
)

// errorDateTimeLayout is the time layout used on job errors
const errorDateTimeLayout = "2006-01-02T15:04:05-07:00"

// OutputGroupType is a custom type for OutputGroup type field values
type OutputGroupType string

const (
	// FileOutputGroupType is the value for the type field on OutputGroup
	// for jobs with a file output
	FileOutputGroupType = OutputGroupType("file_group_settings")
	// AppleLiveOutputGroupType is the value for the type field on OutputGroup
	// for jobs with Apple's HTTP Live Streaming (HLS) output
	AppleLiveOutputGroupType = OutputGroupType("apple_live_group_settings")
)

// Container is the Video container type for a job
type Container string

const (
	// AppleHTTPLiveStreaming is the container for HLS video files
	AppleHTTPLiveStreaming = Container("m3u8")
	// MPEG4 is the container for MPEG-4 video files
	MPEG4 = Container("mp4")
)

// JobErrorDateTime is a custom time struct to be used on Media items
type JobErrorDateTime struct {
	time.Time
}

// MarshalXML implementation on JobErrorDateTime to skip "zero" time values
func (jdt JobErrorDateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if jdt.IsZero() {
		return nil
	}
	e.EncodeElement(jdt.Time, start)
	return nil
}

// UnmarshalXML implementation on JobErrorDateTime to use errorDateTimeLayout
func (jdt *JobErrorDateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var content string
	if err := d.DecodeElement(&content, &start); err != nil {
		return err
	}
	if content == "" {
		jdt.Time = time.Time{}
		return nil
	}
	if content == "0001-01-01T00:00:00Z" {
		jdt.Time = time.Time{}
		return nil
	}
	jdt.Time, err = time.Parse(errorDateTimeLayout, content)
	return err
}

// GetJobs returns a list of the user's jobs
func (c *Client) GetJobs() (*JobList, error) {
	var result *JobList
	err := c.do("GET", "/jobs", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetJob returns metadata on a single job
func (c *Client) GetJob(jobID string) (*Job, error) {
	var result *Job
	err := c.do("GET", "/jobs/"+jobID, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// PostJob sends a single job to the current Elemental
// Cloud deployment for processing
func (c *Client) PostJob(job *Job) (*Job, error) {
	var result *Job
	err := c.do("POST", "/jobs", *job, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetID is a convenience function to parse the job id
// out of the Href attribute in Job
func (j *Job) GetID() string {
	if j.Href != "" {
		hrefData := strings.Split(j.Href, "/")
		if len(hrefData) > 1 {
			return hrefData[len(hrefData)-1]
		}
	}
	return ""
}

// JobList represents the response returned by
// a query for the list of jobs
type JobList struct {
	XMLName xml.Name `xml:"job_list"`
	Empty   string   `xml:"empty,omitempty"`
	Job     []Job    `xml:"job"`
}

// Job represents a job to be sent to Elemental Cloud
type Job struct {
	XMLName         xml.Name         `xml:"job"`
	Href            string           `xml:"href,attr,omitempty"`
	Input           Input            `xml:"input,omitempty"`
	Priority        int              `xml:"priority,omitempty"`
	OutputGroup     OutputGroup      `xml:"output_group,omitempty"`
	StreamAssembly  []StreamAssembly `xml:"stream_assembly,omitempty"`
	Status          string           `xml:"status,omitempty"`
	Submitted       DateTime         `xml:"submitted,omitempty"`
	StartTime       DateTime         `xml:"start_time,omitempty"`
	CompleteTime    DateTime         `xml:"complete_time,omitempty"`
	ErroredTime     DateTime         `xml:"errored_time,omitempty"`
	PercentComplete int              `xml:"pct_complete,omitempty"`
	ErrorMessages   []JobError       `xml:"error_messages,omitempty"`
}

// JobError represents an individual error on a job
type JobError struct {
	Code      int              `xml:"error>code,omitempty"`
	CreatedAt JobErrorDateTime `xml:"error>created_at,omitempty"`
	Message   string           `xml:"error>message,omitempty"`
}

// Input represents the spec for the job's input
type Input struct {
	FileInput Location `xml:"file_input,omitempty"`
}

// Location defines where a file is or needs to be.
// Username and Password are required for certain
// protocols that require authentication, like S3
type Location struct {
	URI      string `xml:"uri,omitempty"`
	Username string `xml:"username,omitempty"`
	Password string `xml:"password,omitempty"`
}

// OutputGroup is a list of the indended outputs for the job
type OutputGroup struct {
	Order                  int                    `xml:"order,omitempty"`
	FileGroupSettings      FileGroupSettings      `xml:"file_group_settings,omitempty"`
	AppleLiveGroupSettings AppleLiveGroupSettings `xml:"apple_live_group_settings,omitempty"`
	Type                   OutputGroupType        `xml:"type,omitempty"`
	Output                 []Output               `xml:"output,omitempty"`
}

// FileGroupSettings define where the file job output should go
type FileGroupSettings struct {
	Destination Location `xml:"destination,omitempty"`
}

// AppleLiveGroupSettings define where the HLS job output should go
type AppleLiveGroupSettings struct {
	Destination     Location `xml:"destination,omitempty"`
	SegmentDuration uint     `xml:"segment_length,omitempty"`
}

// Output defines the different processing stream assemblies
// for the job
type Output struct {
	StreamAssemblyName string    `xml:"stream_assembly_name,omitempty"`
	NameModifier       string    `xml:"name_modifier,omitempty"`
	Order              int       `xml:"order,omitempty"`
	Extension          string    `xml:"extension,omitempty"`
	Container          Container `xml:"container,omitempty"`
}

// StreamAssembly defines how each processing stream should behave
type StreamAssembly struct {
	Name   string `xml:"name,omitempty"`
	Preset string `xml:"preset,omitempty"`
}

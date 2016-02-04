package elementalcloud

// GetJobsResponse represents the response returned by
// a query for job metadata
type GetJobsResponse struct {
	Empty string `xml:"empty,omitempty,omitempty"`
}

// PostJobResponse represents the response returned when
// creating a new job
type PostJobResponse struct {
}

// GetJobs returns a list of the user's jobs
func (c *Client) GetJobs() (*GetJobsResponse, error) {
	var result *GetJobsResponse
	err := c.do("GET", "/jobs", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetJob returns metadata on a single job
func (c *Client) GetJob(jobID string) (*GetJobsResponse, error) {
	var result *GetJobsResponse
	err := c.do("GET", "/jobs/"+jobID, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// PostJob sends a single job to the current Elemental
// Cloud deployment for processing
func (c *Client) PostJob(job Job) (*PostJobResponse, error) {
	var result *PostJobResponse
	err := c.do("POST", "/jobs", job, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Job represents a job to be sent to Elemental Cloud
type Job struct {
	Input          Input            `xml:"input,omitempty"`
	Priority       int              `xml:"priority,omitempty"`
	OutputGroup    OutputGroup      `xml:"output_group,omitempty"`
	StreamAssembly []StreamAssembly `xml:"stream_assembly,omitempty"`
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
	Order             int               `xml:"order,omitempty"`
	FileGroupSettings FileGroupSettings `xml:"file_group_settings,omitempty"`
	Type              string            `xml:"type,omitempty"`
	Output            []Output          `xml:"output,omitempty"`
}

// FileGroupSettings define where the file job output should go
type FileGroupSettings struct {
	Destination Location `xml:"destination,omitempty"`
}

// Output defines the different processing stream assemblies
// for the job
type Output struct {
	StreamAssemblyName string `xml:"stream_assembly_name,omitempty"`
	NameModifier       string `xml:"name_modifier,omitempty"`
	Order              int    `xml:"order,omitempty"`
	Extension          string `xml:"extension,omitempty"`
}

// StreamAssembly defines how each processing stream should behave
type StreamAssembly struct {
	Name   string `xml:"name,omitempty"`
	Preset string `xml:"preset,omitempty"`
}

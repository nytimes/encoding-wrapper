package elementalcloud

// GetJobsResponse represents the response returned by the GET /jobs REST operation.
type GetJobsResponse struct {
	Empty string `xml:"empty,omitempty"`
}

// GetJobs returns a list of the user's jobs in the current Elemental Cloud deployment.
func (c *Client) GetJobs() (*GetJobsResponse, error) {
	var result *GetJobsResponse
	err := c.do("GET", "/jobs", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetJob returns metadata on a single job in the current Elemental Cloud deployment.
func (c *Client) GetJob(jobID string) (*GetJobsResponse, error) {
	var result *GetJobsResponse
	err := c.do("GET", "/jobs/"+jobID, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

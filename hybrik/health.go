package hybrik

// HealthCheck Check to see that the service is responsive
func (c *Client) HealthCheck() error {
	// For now, just call list jobs. If this errors, then we can consider the service unhealthy
	_, err := c.client.CallAPI("GET", "/jobs/info", nil, nil)
	if err != nil {
		return err
	}

	return nil
}

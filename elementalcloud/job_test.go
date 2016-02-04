package elementalcloud

import (
	"net/http"

	"gopkg.in/check.v1"
)

func (s *S) TestGetJobError(c *check.C) {
	errorResponse := `<?xml version="1.0" encoding="UTF-8"?>
<errors>
  <error type="ActiveRecord::RecordNotFound">Couldn't find Job with id=1</error>
</errors>`
	server, _ := s.startServer(http.StatusNotFound, errorResponse)
	defer server.Close()
	client := NewClient(server.URL, "myuser", "secret-key", 45)

	getJobsResponse, err := client.GetJob("1")
	c.Assert(getJobsResponse, check.IsNil)
	c.Assert(err, check.DeepEquals, &APIError{
		Status: http.StatusNotFound,
		Errors: errorResponse,
	})
}

func (s *S) TestGetJobsOnEmptyList(c *check.C) {
	server, _ := s.startServer(http.StatusOK, `<?xml version="1.0" encoding="UTF-8"?>
<job_list>
  <empty>There are currently no jobs</empty>
</job_list>`)
	defer server.Close()
	client := NewClient(server.URL, "myuser", "secret-key", 45)

	getJobsResponse, err := client.GetJobs()
	c.Assert(err, check.IsNil)
	c.Assert(getJobsResponse, check.DeepEquals, &GetJobsResponse{
		Empty: "There are currently no jobs",
	})
}

func (s *S) TestPostJob(c *check.C) {
	server, _ := s.startServer(http.StatusOK, `<response>job added</response>`)
	defer server.Close()

	client := NewClient(server.URL, "myuser", "secret-key", 45)
	jobInput := Job{
		Input: Input{
			FileInput: Location{
				URI:      "http://another.non.existent/video.mp4",
				Username: "user",
				Password: "pass123",
			},
		},
		Priority: 50,
		OutputGroup: OutputGroup{
			Order: 1,
			FileGroupSettings: FileGroupSettings{
				Destination: Location{
					URI:      "http://destination/video.mp4",
					Username: "user",
					Password: "pass123",
				},
			},
			Type: "file_group_settings",
			Output: []Output{
				{
					StreamAssemblyName: "stream_1",
					NameModifier:       "_high",
					Order:              1,
					Extension:          ".mp4",
				},
			},
		},
		StreamAssembly: []StreamAssembly{
			{
				Name:   "stream_1",
				Preset: "17",
			},
		},
	}

	postJobResponse, err := client.PostJob(jobInput)

	c.Assert(err, check.IsNil)
	c.Assert(postJobResponse, check.NotNil)
}

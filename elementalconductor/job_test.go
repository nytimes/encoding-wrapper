package elementalconductor

import (
	"encoding/xml"
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
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

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
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	getJobsResponse, err := client.GetJobs()
	c.Assert(err, check.IsNil)
	c.Assert(getJobsResponse, check.DeepEquals, &JobList{
		XMLName: xml.Name{
			Local: "job_list",
		},
		Empty: "There are currently no jobs",
	})
}

func (s *S) TestPostJob(c *check.C) {
	jobResponseXML := `<job>
    <input>
        <file_input>
            <uri>http://another.non.existent/video.mp4</uri>
            <username>user</username>
            <password>pass123</password>
        </file_input>
    </input>
    <priority>50</priority>
    <output_group>
        <order>1</order>
        <file_group_settings>
            <destination>
                <uri>http://destination/video.mp4</uri>
                <username>user</username>
                <password>pass123</password>
            </destination>
        </file_group_settings>
        <type>file_group_settings</type>
        <output>
            <stream_assembly_name>stream_1</stream_assembly_name>
            <name_modifier>_high</name_modifier>
            <order>1</order>
            <extension>.mp4</extension>
        </output>
    </output_group>
    <stream_assembly>
        <name>stream_1</name>
        <preset>17</preset>
    </stream_assembly>
</job>`
	server, _ := s.startServer(http.StatusOK, jobResponseXML)
	defer server.Close()
	jobInput := Job{
		XMLName: xml.Name{
			Local: "job",
		},
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
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	postJobResponse, err := client.PostJob(&jobInput)
	c.Assert(err, check.IsNil)
	c.Assert(postJobResponse, check.NotNil)
	c.Assert(postJobResponse, check.DeepEquals, &jobInput)
}

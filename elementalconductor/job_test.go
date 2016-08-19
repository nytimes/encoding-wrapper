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

func (s *S) TestCreateJob(c *check.C) {
	jobResponseXML := `<job href="/jobs/1">
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
        <apple_live_group_settings>
            <destination>
                <uri>http://destination/video.mp4</uri>
                <username>user</username>
                <password>pass123</password>
            </destination>
        </apple_live_group_settings>
        <type>apple_live_group_settings</type>
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
	server, _ := s.startServer(http.StatusCreated, jobResponseXML)
	defer server.Close()
	jobInput := Job{
		XMLName: xml.Name{
			Local: "job",
		},
		Href: "/jobs/1",
		Input: Input{
			FileInput: Location{
				URI:      "http://another.non.existent/video.mp4",
				Username: "user",
				Password: "pass123",
			},
		},
		Priority: 50,
		OutputGroup: []OutputGroup{
			{
				Order: 1,
				FileGroupSettings: &FileGroupSettings{
					Destination: &Location{
						URI:      "http://destination/video.mp4",
						Username: "user",
						Password: "pass123",
					},
				},
				AppleLiveGroupSettings: &AppleLiveGroupSettings{
					Destination: &Location{
						URI:      "http://destination/video.mp4",
						Username: "user",
						Password: "pass123",
					},
				},
				Type: AppleLiveOutputGroupType,
				Output: []Output{
					{
						StreamAssemblyName: "stream_1",
						NameModifier:       "_high",
						Order:              1,
						Extension:          ".mp4",
					},
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

	postJobResponse, err := client.CreateJob(&jobInput)
	c.Assert(err, check.IsNil)
	c.Assert(postJobResponse, check.NotNil)
	c.Assert(postJobResponse, check.DeepEquals, &jobInput)
}

func (s *S) TestGetJob(c *check.C) {
	jobResponseXML := `<job href="/jobs/1">
    <input>
        <file_input>
            <uri>http://another.non.existent/video.mp4</uri>
            <username>user</username>
            <password>pass123</password>
        </file_input>
    </input>
    <content_duration>
        <input_duration>716</input_duration>
        <clipped_input_duration>716</clipped_input_duration>
        <stream_count>1</stream_count>
        <total_stream_duration>716</total_stream_duration>
        <package_count>1</package_count>
        <total_package_duration>716</total_package_duration>
    </content_duration>
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
	expectedJob := Job{
		XMLName: xml.Name{
			Local: "job",
		},
		Href: "/jobs/1",
		Input: Input{
			FileInput: Location{
				URI:      "http://another.non.existent/video.mp4",
				Username: "user",
				Password: "pass123",
			},
		},
		ContentDuration: ContentDuration{
			InputDuration: 716,
		},
		Priority: 50,
		OutputGroup: []OutputGroup{
			{
				Order: 1,
				FileGroupSettings: &FileGroupSettings{
					Destination: &Location{
						URI:      "http://destination/video.mp4",
						Username: "user",
						Password: "pass123",
					},
				},
				Type: FileOutputGroupType,
				Output: []Output{
					{
						StreamAssemblyName: "stream_1",
						NameModifier:       "_high",
						Order:              1,
						Extension:          ".mp4",
					},
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

	getJobsResponse, err := client.GetJob("1")
	c.Assert(err, check.IsNil)
	c.Assert(getJobsResponse, check.NotNil)
	c.Assert(getJobsResponse, check.DeepEquals, &expectedJob)
}

func (s *S) TestCancelJob(c *check.C) {
	jobResponseXML := `<job href="/jobs/1">
    <status>canceled</status>
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
	server, reqs := s.startServer(http.StatusOK, jobResponseXML)
	defer server.Close()
	expectedJob := Job{
		XMLName: xml.Name{
			Local: "job",
		},
		Href: "/jobs/1",
		Input: Input{
			FileInput: Location{
				URI:      "http://another.non.existent/video.mp4",
				Username: "user",
				Password: "pass123",
			},
		},
		Priority: 50,
		OutputGroup: []OutputGroup{
			{
				Order: 1,
				FileGroupSettings: &FileGroupSettings{
					Destination: &Location{
						URI:      "http://destination/video.mp4",
						Username: "user",
						Password: "pass123",
					},
				},
				Type: FileOutputGroupType,
				Output: []Output{
					{
						StreamAssemblyName: "stream_1",
						NameModifier:       "_high",
						Order:              1,
						Extension:          ".mp4",
					},
				},
			},
		},
		StreamAssembly: []StreamAssembly{
			{
				Name:   "stream_1",
				Preset: "17",
			},
		},
		Status: "canceled",
	}
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	job, err := client.CancelJob("1")
	c.Assert(err, check.IsNil)
	c.Assert(job, check.NotNil)
	c.Assert(job, check.DeepEquals, &expectedJob)

	req := <-reqs
	c.Assert(req.req.Method, check.Equals, "POST")
	c.Assert(req.req.URL.Path, check.Equals, "/api/jobs/1/cancel")
	c.Assert(string(req.body), check.Equals, "<cancel></cancel>")
}

func (s *S) TestCancelJobError(c *check.C) {
	errorResponse := `<?xml version="1.0" encoding="UTF-8"?>
<errors>
  <error type="ActiveRecord::RecordNotFound">Couldn't find Job with id=1</error>
</errors>`
	server, _ := s.startServer(http.StatusNotFound, errorResponse)
	defer server.Close()
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	job, err := client.CancelJob("1")
	c.Assert(job, check.IsNil)
	c.Assert(err, check.DeepEquals, &APIError{
		Status: http.StatusNotFound,
		Errors: errorResponse,
	})
}

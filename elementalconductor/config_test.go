package elementalconductor

import (
	"encoding/xml"
	"net/http"
	"reflect"
	"testing"
)

func TestGetCloudConfig(t *testing.T) {
	data := `<?xml version="1.0" encoding="UTF-8"?>
<cloud_config>
  <authorized_node_count>500</authorized_node_count>
  <max_cluster_size>30</max_cluster_size>
  <min_cluster_size>4</min_cluster_size>
  <worker_variant>production_server_cloud</worker_variant>
</cloud_config>`
	server, requests := startServer(http.StatusOK, data)
	defer server.Close()
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")
	config, err := client.GetCloudConfig()
	if err != nil {
		t.Fatal(err)
	}
	expectedConfig := CloudConfig{
		XMLName:             xml.Name{Local: "cloud_config"},
		AuthorizedNodeCount: 500,
		MaxNodes:            30,
		MinNodes:            4,
		WorkerVariant:       "production_server_cloud",
	}
	if !reflect.DeepEqual(*config, expectedConfig) {
		t.Errorf("wrong config returned\nwant %#v\ngot  %#v", expectedConfig, *config)
	}

	fakeReq := <-requests
	if fakeReq.req.Method != http.MethodGet {
		t.Errorf("wrong http method\nwant %q\ngot  %q", http.MethodGet, fakeReq.req.Method)
	}
	if expectedPath := "/api/config/cloud"; fakeReq.req.URL.Path != expectedPath {
		t.Errorf("wrong request path\nwant %q\ngot  %q", expectedPath, fakeReq.req.URL.Path)
	}
}

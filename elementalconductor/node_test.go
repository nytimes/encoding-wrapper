package elementalconductor

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGetNodes(t *testing.T) {
	server, requests := startServer(http.StatusOK, `<?xml version="1.0" encoding="UTF-8"?>
<node_list>
  <node href="/nodes/1">
    <name>Conductor</name>
    <hostname>ip-192-168-1-140</hostname>
    <ip_addr>192.168.1.140</ip_addr>
    <eth0_mac>0E:C8:60:FA:3C:01</eth0_mac>
    <status>active</status>
    <product>Conductor File</product>
    <version>1.0.0v123</version>
    <platform>cloud</platform>
    <packages>
      <package>Audio Normalization Package</package>
      <package>Audio Package</package>
    </packages>
    <licenses>
      <license>srs_norm</license>
      <license>dolby_ac3</license>
      <license>dolby_eac3</license>
      <license>dts</license>
    </licenses>
    <created_at>2015-12-02 15:15:59 -0300</created_at>
  </node>
  <node href="/nodes/31">
    <name>Node 1</name>
    <hostname>ip-192-168-1-141</hostname>
    <ip_addr>192.168.1.141</ip_addr>
    <public_ip_addr>50.10.10.199</public_ip_addr>
    <eth0_mac>0E:C8:60:FA:3C:02</eth0_mac>
    <status>active</status>
    <product>Server</product>
    <version>1.0.0v123</version>
    <platform>cloud</platform>
    <packages>
      <package>Audio Normalization Package</package>
      <package>Audio Package</package>
    </packages>
    <licenses>
      <license>srs_norm</license>
      <license>dolby_ac3</license>
      <license>dolby_eac3</license>
      <license>dts</license>
    </licenses>
    <created_at>2016-03-01 09:42:23 -0300</created_at>
    <backup_groups>
    </backup_groups>
    <running_count>80</running_count>
  </node>
  <node href="/nodes/40">
    <name>Node 2</name>
    <hostname>ip-192-168-1-142</hostname>
    <ip_addr>192.168.1.142</ip_addr>
    <eth0_mac>0E:C8:61:FA:3C:03</eth0_mac>
    <status>active</status>
    <product>Server</product>
    <version>1.0.0v123</version>
    <platform>cloud</platform>
    <packages>
      <package>Audio Normalization Package</package>
      <package>Audio Package</package>
    </packages>
    <licenses>
      <license>srs_norm</license>
      <license>dolby_ac3</license>
      <license>dolby_eac3</license>
      <license>dts</license>
    </licenses>
    <created_at>2016-03-25 21:49:01 -0300</created_at>
    <backup_groups>
    </backup_groups>
    <running_count>120</running_count>
  </node>
</node_list>`)
	defer server.Close()
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")
	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
	}
	expectedNodes := []Node{
		{
			Href:      "/nodes/1",
			Name:      "Conductor",
			HostName:  "ip-192-168-1-140",
			IPAddress: "192.168.1.140",
			Eth0Mac:   "0E:C8:60:FA:3C:01",
			Status:    "active",
			Product:   ProductConductorFile,
			Version:   "1.0.0v123",
			Platform:  "cloud",
			Packages:  []string{"Audio Normalization Package", "Audio Package"},
			Licenses:  []string{"srs_norm", "dolby_ac3", "dolby_eac3", "dts"},
			CreatedAt: DateTime{Time: time.Date(2015, time.December, 2, 18, 15, 59, 0, time.UTC)},
		},
		{
			Href:            "/nodes/31",
			Name:            "Node 1",
			HostName:        "ip-192-168-1-141",
			IPAddress:       "192.168.1.141",
			PublicIPAddress: "50.10.10.199",
			Eth0Mac:         "0E:C8:60:FA:3C:02",
			Status:          "active",
			Product:         ProductServer,
			Version:         "1.0.0v123",
			Platform:        "cloud",
			Packages:        []string{"Audio Normalization Package", "Audio Package"},
			Licenses:        []string{"srs_norm", "dolby_ac3", "dolby_eac3", "dts"},
			CreatedAt:       DateTime{Time: time.Date(2016, time.March, 1, 12, 42, 23, 0, time.UTC)},
			RunningCount:    80,
		},
		{
			Href:         "/nodes/40",
			Name:         "Node 2",
			HostName:     "ip-192-168-1-142",
			IPAddress:    "192.168.1.142",
			Eth0Mac:      "0E:C8:61:FA:3C:03",
			Status:       "active",
			Product:      ProductServer,
			Version:      "1.0.0v123",
			Platform:     "cloud",
			Packages:     []string{"Audio Normalization Package", "Audio Package"},
			Licenses:     []string{"srs_norm", "dolby_ac3", "dolby_eac3", "dts"},
			CreatedAt:    DateTime{Time: time.Date(2016, time.March, 26, 0, 49, 1, 0, time.UTC)},
			RunningCount: 120,
		},
	}
	if !reflect.DeepEqual(nodes, expectedNodes) {
		t.Errorf("wrong nodes returned\nwant %#v\ngot  %#v", expectedNodes, nodes)
	}

	fakeReq := <-requests
	if expectedMethod := "GET"; fakeReq.req.Method != expectedMethod {
		t.Errorf("wrong method used\nwant %q\ngot  %q", expectedMethod, fakeReq.req.Method)
	}
	if expectedPath := "/api/nodes"; fakeReq.req.URL.Path != expectedPath {
		t.Errorf("wrong path used\nwant %q\ngot  %q", expectedPath, fakeReq.req.URL.Path)
	}
}

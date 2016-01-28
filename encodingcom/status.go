package encodingcom

import (
	"errors"
	"strings"
	"time"
)

// StatusResponse is the result of the GetStatus method.
//
// See http://goo.gl/NDsN8h for more details.
type StatusResponse struct {
	MediaID             string
	UserID              string
	SourceFile          string
	MediaStatus         string
	PreviousMediaStatus string
	NotifyURL           string
	CreateDate          time.Time
	StartDate           time.Time
	FinishDate          time.Time
	DownloadDate        time.Time
	UploadDate          time.Time
	TimeLeft            string
	Progress            float64
	TimeLeftCurrentJob  string
	ProgressCurrentJob  float64
	Formats             []FormatStatus
}

// FormatStatus is the status of each formatting input for a given MediaID.
//
// It is part of the StatusResponse type.
type FormatStatus struct {
	ID            string
	Status        string
	CreateDate    time.Time
	StartDate     time.Time
	FinishDate    time.Time
	S3Destination string
	CFDestination string
	Destinations  []DestinationStatus
}

// DestinationStatus represents the status of a given destination.
type DestinationStatus struct {
	Name   string
	Status string
}

// GetStatus returns the status of the given media ids, it returns an slice of
// StatusResponse, the size of the result slice matches the size of input
// slice.
func (c *Client) GetStatus(mediaIDs []string) ([]StatusResponse, error) {
	if len(mediaIDs) == 0 {
		return nil, errors.New("please provide at least one media id")
	}
	var m map[string]map[string][]statusJSON
	err := c.do(&request{
		Action:   "GetStatus",
		MediaID:  strings.Join(mediaIDs, ","),
		Extended: true,
	}, &m)
	if err != nil {
		return nil, err
	}
	apiStatus := m["response"]["job"]
	statusResponse := make([]StatusResponse, len(apiStatus))
	for i, status := range apiStatus {
		statusResponse[i] = status.toStruct()
	}
	return statusResponse, nil
}

type statusJSON struct {
	MediaID             string             `json:"id"`
	UserID              string             `json:"userid"`
	SourceFile          string             `json:"sourcefile"`
	MediaStatus         string             `json:"status"`
	PreviousMediaStatus string             `json:"prevstatus"`
	NotifyURL           string             `json:"notifyurl"`
	CreateDate          MediaDateTime      `json:"created"`
	StartDate           MediaDateTime      `json:"started"`
	FinishDate          MediaDateTime      `json:"finished"`
	DownloadDate        MediaDateTime      `json:"downloaded"`
	UploadDate          MediaDateTime      `json:"uploaded"`
	TimeLeft            string             `json:"time_left"`
	Progress            float64            `json:"progress,string"`
	TimeLeftCurrentJob  string             `json:"time_left_current"`
	ProgressCurrentJob  float64            `json:"progress_current,string"`
	Formats             []formatStatusJSON `json:"format"`
}

func (s *statusJSON) toStruct() StatusResponse {
	resp := StatusResponse{
		MediaID:             s.MediaID,
		UserID:              s.UserID,
		SourceFile:          s.SourceFile,
		MediaStatus:         s.MediaStatus,
		PreviousMediaStatus: s.PreviousMediaStatus,
		NotifyURL:           s.NotifyURL,
		CreateDate:          s.CreateDate.Time,
		StartDate:           s.StartDate.Time,
		FinishDate:          s.FinishDate.Time,
		DownloadDate:        s.DownloadDate.Time,
		UploadDate:          s.UploadDate.Time,
		TimeLeft:            s.TimeLeft,
		Progress:            s.Progress,
		TimeLeftCurrentJob:  s.TimeLeftCurrentJob,
		ProgressCurrentJob:  s.ProgressCurrentJob,
		Formats:             make([]FormatStatus, len(s.Formats)),
	}
	for i, formatStatus := range s.Formats {
		format := FormatStatus{
			ID:            formatStatus.ID,
			Status:        formatStatus.Status,
			CreateDate:    formatStatus.CreateDate.Time,
			StartDate:     formatStatus.StartDate.Time,
			FinishDate:    formatStatus.FinishDate.Time,
			S3Destination: formatStatus.S3Destination,
			CFDestination: formatStatus.CFDestination,
			Destinations:  make([]DestinationStatus, len(formatStatus.Destinations)),
		}
		for i, dest := range formatStatus.Destinations {
			destStatus := DestinationStatus{Name: dest}
			if i < len(formatStatus.DestinationsStatus) {
				destStatus.Status = formatStatus.DestinationsStatus[i]
			}
			format.Destinations[i] = destStatus
		}
		resp.Formats[i] = format
	}
	return resp
}

type formatStatusJSON struct {
	ID                 string        `json:"id"`
	Status             string        `json:"status"`
	CreateDate         MediaDateTime `json:"created"`
	StartDate          MediaDateTime `json:"started"`
	FinishDate         MediaDateTime `json:"finished"`
	S3Destination      string        `json:"s3_destination"`
	CFDestination      string        `json:"cf_destination"`
	Destinations       []string      `json:"destination"`
	DestinationsStatus []string      `json:"destination_status"`
}

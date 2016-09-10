package elementalconductor

import (
	"encoding/xml"
	"time"
)

const dateTimeLayout = "2006-01-02 15:04:05 -0700"

// DateTime is a custom struct for representing time within ElementalConductor.
// It customizes marshalling, and always store the underlying time in UTC.
type DateTime struct {
	time.Time
}

// MarshalXML implementation on JobDateTime to skip "zero" time values
func (jdt DateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if jdt.IsZero() {
		return nil
	}
	e.EncodeElement(jdt.Time, start)
	return nil
}

// UnmarshalXML implementation on JobDateTime to use dateTimeLayout
func (jdt *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var content string
	err := d.DecodeElement(&content, &start)
	if err != nil {
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
	jdt.Time, err = time.Parse(dateTimeLayout, content)
	if err != nil {
		return err
	}
	jdt.Time = jdt.Time.UTC()
	return nil
}

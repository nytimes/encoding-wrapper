package elementalconductor

import (
	"encoding/xml"
	"fmt"
	"time"

	"gopkg.in/check.v1"
)

func (s *S) TestDateTimeUnmarshalXML(c *check.C) {
	var tests = []struct {
		input    string
		expected time.Time
	}{
		{
			"2016-02-01 11:59:20 -0800",
			time.Date(2016, time.February, 1, 19, 59, 20, 0, time.UTC),
		},
		{
			"2016-02-01 00:25:00 +0300",
			time.Date(2016, time.January, 31, 21, 25, 0, 0, time.UTC),
		},
	}
	for _, test := range tests {
		var output struct {
			XMLName xml.Name `xml:"item"`
			Date    DateTime `xml:"date"`
		}
		input := fmt.Sprintf("<item><date>%s</date></item>", test.input)
		err := xml.Unmarshal([]byte(input), &output)
		c.Check(err, check.IsNil)
		c.Check(output.Date.Time, check.DeepEquals, test.expected)
	}
}

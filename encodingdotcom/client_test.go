package encodingdotcom

import (
	"encoding/json"
	"testing"

	"gopkg.in/check.v1"
)

type S struct{}

func Test(t *testing.T) {
	check.TestingT(t)
}

var _ = check.Suite(&S{})

func (s *S) TestYesNoBoolean(c *check.C) {
	bTrue := YesNoBoolean(true)
	bFalse := YesNoBoolean(false)
	data, err := json.Marshal(bTrue)
	c.Assert(err, check.IsNil)
	c.Assert(string(data), check.Equals, `"yes"`)
	data, err = json.Marshal(bFalse)
	c.Assert(err, check.IsNil)
	c.Assert(string(data), check.Equals, `"no"`)
}

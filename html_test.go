package main

import (
	"testing"

	"github.com/mgutz/ansi"
	"github.com/stretchr/testify/suite"
)

type HTMLTestSuite struct {
	suite.Suite
}

func TestHTMLTestSuite(t *testing.T) {
	suite.Run(t, new(HTMLTestSuite))
}

func (suite *HTMLTestSuite) TestHTMLtoText() {
	cases := map[string]string{
		`<pre><code>this is some code</code></pre>`: ansi.Color("this is some code", "gray"),
		`<a href="example.com">example</a>`:         ansi.Color("example", "blue"),

		`<p>Thanks <spark-mention data-object-type="person" data-object-id="FAKE_PERSON_ID">Kevin</spark-mention></p>`: "Thanks " + ansi.Color("this is some code", "gray"),
	}

	for input, expected := range cases {
		actual := HTMLtoText(input)
		suite.Equal(expected, actual)
	}
}

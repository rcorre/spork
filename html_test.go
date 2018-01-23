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
	conf := &Config{
		MessageFormat: messageFormat{
			emph:    "white+h",
			strong:  "white+b",
			code:    "white:gray",
			quote:   "gray",
			link:    "blue",
			mention: "blue+b",
		},
	}

	cases := map[string]string{
		`one<br/>two`: "one\ntwo",

		`<code>this is some code</code>`:    ansi.Color("this is some code", "white:gray"),
		`<a href="example.com">example</a>`: ansi.Color("example", "blue"),

		"<pre><code>x=5\ny=257</code></pre>": ansi.Color("x=5  \ny=257", "white:gray"),

		`<p>Thanks <spark-mention data-object-type="person" data-object-id="FAKE_PERSON_ID">Kevin</spark-mention></p>`: "Thanks " + ansi.Color("Kevin", "blue+b"),
	}

	for input, expected := range cases {
		actual := HTMLtoText(input, conf)
		suite.Equal(expected, actual)
	}
}

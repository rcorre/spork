package main

import (
	"testing"

	"github.com/mgutz/ansi"
	"github.com/stretchr/testify/suite"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

type MarkdownTestSuite struct {
	suite.Suite
}

func TestMarkdownTestSuite(t *testing.T) {
	suite.Run(t, new(MarkdownTestSuite))
}

var cases = map[string]string{
	"foo":     "foo",
	"*foo*":   ansi.Color("white+h", "foo"),
	"_foo_":   ansi.Color("white+h", "foo"),
	"**foo**": ansi.Color("white+b", "foo"),
}

func (suite *MarkdownTestSuite) TestRender() {
	md := NewMarkdown(nil)
	input := "foo"
	expected := "foo"
	actual := blackfriday.Run([]byte(input), blackfriday.WithRenderer(md))
	suite.Equal(expected, string(actual))
}

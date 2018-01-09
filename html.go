package main

import (
	"strings"

	"github.com/romana/rlog"
	"golang.org/x/net/html"
)

func HTMLtoText(in string) string {
	rlog.Debugf("Parsing HTML %s", in)
	r := strings.NewReader(in)
	doc, err := html.Parse(r)
	if err != nil {
		rlog.Errorf("Failed to parse html %s: %v", in, err)
		return ""
	}

	var s string
	var f func(*html.Node)
	f = func(n *html.Node) {
		rlog.Debugf("\tHandling HTML node %v", n)

		if n.Type == html.ElementNode && n.Data == "a" {
			// Do something with n...
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return s
}

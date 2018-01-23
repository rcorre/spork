package main

import (
	"strings"

	"github.com/mgutz/ansi"
	"github.com/romana/rlog"
	"golang.org/x/net/html"
)

type attributes struct {
	code    bool
	pre     bool
	link    bool
	mention bool
}

func HTMLtoText(in string, conf *Config) string {
	rlog.Debugf("Parsing HTML %s", in)
	r := strings.NewReader(in)
	doc, err := html.Parse(r)
	if err != nil {
		rlog.Errorf("Failed to parse html %s: %v", in, err)
		return ""
	}

	var s string
	var f func(*html.Node, attributes)
	f = func(n *html.Node, attr attributes) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "code":
				attr.code = true
			case "pre":
				attr.pre = true
			case "a":
				attr.link = true
			case "spark-mention":
				attr.mention = true
			case "br":
				s += "\n"
			}
		} else if n.Type == html.TextNode {
			if attr.pre {
				// pad all lines to the same len to create a block of background color
				lines := strings.Split(n.Data, "\n")

				maxlen := 0
				for _, l := range lines {
					if n := len(l); n > maxlen {
						maxlen = n
					}
				}

				for i, l := range lines {
					n := maxlen - len(l)
					lines[i] = l + strings.Repeat(" ", n)
				}

				res := strings.Join(lines, "\n")
				s += ansi.Color(res, conf.MessageFormat.code)
			} else if attr.code {
				s += ansi.Color(n.Data, conf.MessageFormat.code)
			} else if attr.link {
				s += ansi.Color(n.Data, conf.MessageFormat.link)
			} else if attr.mention {
				s += ansi.Color(n.Data, conf.MessageFormat.mention)
			} else {
				s += n.Data
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, attr)
		}
	}
	f(doc, attributes{})
	return s
}

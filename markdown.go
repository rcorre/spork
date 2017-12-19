package main

import (
	"io"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

type markdown struct {
	typ blackfriday.NodeType
}

func NewMarkdown(conf *Config) blackfriday.Renderer {
	return &markdown{}
}

func (m *markdown) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.Document:
	case blackfriday.BlockQuote:
	case blackfriday.List:
	case blackfriday.Item:
	case blackfriday.Paragraph:
	case blackfriday.Heading:
	case blackfriday.HorizontalRule:
	case blackfriday.Emph:
	case blackfriday.Strong:
	case blackfriday.Del:
	case blackfriday.Link:
	case blackfriday.Image:
	case blackfriday.Text:
		w.Write(node.Literal)
	case blackfriday.HTMLBlock:
	case blackfriday.CodeBlock:
	case blackfriday.Softbreak:
	case blackfriday.Hardbreak:
	case blackfriday.Code:
	case blackfriday.HTMLSpan:
	case blackfriday.Table:
	case blackfriday.TableCell:
	case blackfriday.TableHead:
	case blackfriday.TableBody:
	case blackfriday.TableRow:
	}

	return blackfriday.GoToNext
}

func (m *markdown) RenderHeader(w io.Writer, ast *blackfriday.Node) {
}

func (m *markdown) RenderFooter(w io.Writer, ast *blackfriday.Node) {
}

/*
	Extract first paragraph from markdown as summary
*/

package utils

import (
	"io"

	"github.com/russross/blackfriday"
)

// SummaryRender is the rendering interface
type SummaryRender struct {
}

// RenderNode extract first paragraph
func (r SummaryRender) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	if node.Parent != nil && node.Parent.Type == blackfriday.Paragraph {
		w.Write(node.Literal)
		return blackfriday.GoToNext
	}
	if node.Type == blackfriday.Paragraph && entering == false {
		return blackfriday.Terminate
	}
	return blackfriday.GoToNext
}

// RenderHeader can produce extra content before main document
func (r SummaryRender) RenderHeader(w io.Writer, ast *blackfriday.Node) {
}

// RenderFooter is a symmetric counterpart of RenderHeader.
func (r SummaryRender) RenderFooter(w io.Writer, ast *blackfriday.Node) {
}

/*
	Extract first paragraph from markdown as summary
*/

package utils

import (
	"fmt"
	"io"

	"github.com/russross/blackfriday"
)

// SummaryRenderer is the rendering interface
type SummaryRenderer struct {
}

// RenderNode extract first paragraph
func (r *SummaryRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	if node.Parent != nil && node.Parent.Type == blackfriday.Paragraph {
		switch node.Type {
		case blackfriday.Text:
			w.Write(node.Literal)
		case blackfriday.Code:
			w.Write([]byte(fmt.Sprintf("<code>%s</code>", node.Literal)))
		}
		return blackfriday.GoToNext
	}
	if node.Type == blackfriday.Paragraph && entering == false {
		return blackfriday.Terminate
	}
	return blackfriday.GoToNext
}

// RenderHeader can produce extra content before main document
func (r *SummaryRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
}

// RenderFooter is a symmetric counterpart of RenderHeader.
func (r *SummaryRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
}

// NewSummaryRenderer creates and configures an HTMLRenderer object, which satisfies the Renderer interface.
func NewSummaryRenderer() *SummaryRenderer {
	return &SummaryRenderer{}
}

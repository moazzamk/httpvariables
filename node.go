package variable

import (
	"github.com/pytlesk4/m"
	"fmt"
)

type NodeType int


const (
	SingleAction  NodeType = iota
	DoubleAction  NodeType = iota
	BlocAction    NodeType = iota
	NormalNode    NodeType = iota
	BlockTextNode NodeType = iota
)

// Node interface
type Node interface {
	String() string
}

type ActionNodeOption func(node *ActionNode)

//Action node. This node is used to store a single or double action part of the template
type ActionNode struct {
	Type     NodeType
	Text     string
	Dict    map[string]interface{}
}

// Returns a new action node with the specified options
// By default it returns a single action node
func NewActionNode(options ...ActionNodeOption) *ActionNode {
	ret := &ActionNode{
		Type: SingleAction,
	}

	for _, fn := range options {
		fn(ret)
	}

	return ret
}

// WithActionType is used as an argument for NewActionNode. It helps create a new action node with
// the specified type
func WithActionType(nodeType NodeType) ActionNodeOption {
	return func (node *ActionNode) {
		node.Type = nodeType
	}
}

// WithActionText is used as an argument for NewActionNode. It helps create a new action node with
// the specified text
func WithActionText(nodeText string) ActionNodeOption {
	return func (node *ActionNode) {
		node.Text = nodeText
	}
}

// WithActionType is used as an argument for NewActionNode. It helps create a new action node with
// the specified dictionary (of template variables)
func WithActionDictionary(dict map[string]interface{}) ActionNodeOption {
	return func (node *ActionNode) {
		node.Dict = dict
	}
}

// String returns rendered template as a string
func (r *ActionNode) String() string {
	val := m.Get(r.Dict, r.Text)
	if val == nil {
		return `{` + r.Text + `}`
	}

	if r.Type == DoubleAction {
		val = m.Get(r.Dict, fmt.Sprintf("%v", val))
		if val == nil {
			return `{{` + r.Text + `}`
		}
	}

	return fmt.Sprintf("%v", val)
}

type TextNodeOption func(node *TextNode)

// A regular text node
type TextNode struct {
	Text string
}

// NewTextNode returns a new text node
func NewTextNode(options ...TextNodeOption) *TextNode {
	ret := &TextNode{}

	for _, fn := range options {
		fn(ret)
	}

	return ret
}

// WithText is used as an argument for NewTextNode. It helps create a new text node with
// the specified text
func WithText(text string) TextNodeOption {
	return func (node *TextNode) {
		node.Text = text
	}
}

// String returns rendered template as a string
func (r *TextNode) String() string {
	return r.Text
}



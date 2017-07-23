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

func WithActionType(nodeType NodeType) ActionNodeOption {
	return func (node *ActionNode) {
		node.Type = nodeType
	}
}

func WithActionText(nodeText string) ActionNodeOption {
	return func (node *ActionNode) {
		node.Text = nodeText
	}
}

func WithActionDictionary(dict map[string]interface{}) ActionNodeOption {
	return func (node *ActionNode) {
		node.Dict = dict
	}
}

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

func NewTextNode(options ...TextNodeOption) *TextNode {
	ret := &TextNode{}

	for _, fn := range options {
		fn(ret)
	}

	return ret
}

func WithText(text string) TextNodeOption {
	return func (node *TextNode) {
		node.Text = text
	}
}

func (r *TextNode) String() string {
	return r.Text
}


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

// WithActionDictionary is used as an argument for NewActionNode. It helps create a new action node with
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
	bytes []byte
}


// NewTextNode returns a new instance of a TextNode
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
		node.bytes = []byte(text)
	}
}

// WithBytes returns an option for TextNode when creating a new text node
func WithBytes(bytes []byte) TextNodeOption {
	return func (node *TextNode) {
		node.bytes = bytes
	}
}

// String returns rendered template as a string
func (r *TextNode) String() string {
	return string(r.bytes)
}

func (r *TextNode) Bytes() []byte {
	return r.bytes
}

// BlockExpressionEvaluator is a function that evaluates the expression given to BlockNode
type BlockExpressionEvaluator func (string) bool
type BlockNodeOption func (node *BlockNode)
type BlockNode struct {
	evaluator BlockExpressionEvaluator
	expression string
	children []Node
	dict map[string]interface{}
}

// NewBlockNode returns a new instance of BlockNode
func NewBlockNode(options ...BlockNodeOption) *BlockNode {
	ret := &BlockNode{}
	for _, val := range options {
		val(ret)
	}

	return ret
}

// WithDict is used as an argument for NewBlock. It helps create a new block node with
// the specified dictionary (of template variables)
func WithDict(dict map[string]interface{}) BlockNodeOption {
	return func(node *BlockNode) {
		node.dict = dict
	}
}

// WithDict is used as an argument for NewBlock. It helps create a new block node with
// the specified expression (of template variables)
func WithExpression(expression string) BlockNodeOption {
	return func(node *BlockNode) {
		node.expression = expression
	}
}

// Expression() gets the variable that needs to exist in data dictionary
// for block node to render child block
func (r *BlockNode) Expression() string {
	return r.expression
}

// SetExpression() sets the variable that needs to exist in data dictionary
// for block node to render child block
func (r *BlockNode) SetExpression(val string) *BlockNode {
	r.expression = val

	return r
}

// AddChild adds a child node to BlockNode
func (r *BlockNode) AddChild(node Node) {
	r.children = append(r.children, node)
}

// String returns rendered template of block node. If expression is blank string then rendered block is returned.
// If expression for this BlockNode is not empty string then the variable in the expression must exist in data dictionary
// for the block to be rendered
func (r *BlockNode) String() string {
	if _, ok := r.dict[r.expression]; r.expression != `` && !ok {
		return ``
	}

	ret := ``
	for _, val := range r.children {
		ret += val.String()
	}

	return ret
}

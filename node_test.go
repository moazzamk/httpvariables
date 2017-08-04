package variable

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_ActionNode_ignores_template_holes_with_no_variables(t *testing.T) {
	mappy := make(map[string]interface{})

	node := NewActionNode(WithActionDictionary(mappy),
						WithActionText("yo"))

	assert.Equal(t, "{yo}", node.String())
}

func Test_ActionNode_variable_of_variable_syntax(t *testing.T) {
	mappy := make(map[string]interface{})
	mappy["a"] = "b"
	mappy["b"] = "c"

	node := NewActionNode(WithActionDictionary(mappy),
						WithActionType(DoubleAction),
						WithActionText("a"))

	assert.Equal(t, "c", node.String())
}

func Test_BlockNode_is_initializable(t *testing.T) {
	assert.IsType(t, &BlockNode{}, NewBlockNode())
}

func Test_BlockNode_renders_children(t *testing.T) {
	node := NewTextNode(WithText("hello"))
	parent := NewBlockNode()

	parent.AddChild(node)

	assert.Equal(t, "hello", parent.String())
}

func Test_BlockNode_wont_render_children_if_expression_doesnt_return_true(t *testing.T) {
	dict := make(map[string]interface{})
	dict["hi"] = "hello"

	node := NewBlockNode(WithDict(dict),
						WithExpression("yo"))
	parent := NewBlockNode()

	parent.AddChild(node)

	assert.Equal(t, "", parent.String())

}

func Test_BlockNode_renders_children_if_expression_returns_true(t *testing.T) {
	dict := make(map[string]interface{})
	dict["hi"] = "hello"

	node := NewTextNode(WithText("jojo"))
	parent := NewBlockNode(WithDict(dict),
							WithExpression("hi"))

	parent.AddChild(node)

	assert.Equal(t, "jojo", parent.String())
}



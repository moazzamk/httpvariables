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

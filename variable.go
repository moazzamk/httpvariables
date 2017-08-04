package variable

import (
	"strings"
	"bufio"
	"fmt"
	"net/http"

	"encoding/json"
	"bytes"
	"io/ioutil"
)

const (
	Brace = byte('{')
	Hash = byte('#')
	FSlash = byte('/')

)

// Replaces template variables in a request.
// Check ReplaceVariable() to see how templates get populated
func PopulateRequestTemplate(req *http.Request, variables string) {
	if variables == `` {
		return
	}

	variableData := make(map[string]interface{})
	err := json.Unmarshal([]byte(variables), &variableData)
	if err != nil {
		fmt.Println("JSON MARSHAL ERROR ", err)
		//fmt.Println("EMPTTYTYTY2")
		return
	}


	if req.Body != nil {
		body, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewReader([]byte(ReplaceVariable(string(body), variableData))))
	}


	req.Method = ReplaceVariable(req.Method, variableData)
	req.URL.Host = ReplaceVariable(req.URL.Host, variableData)
	req.URL.RawQuery = ReplaceVariable(req.URL.RawQuery, variableData)
	req.URL.Path = ReplaceVariable(req.URL.Path, variableData)

}

// Replaces templates variables in a string
// Examples:
// {key} - will be replaced with a template variable called key
// {{key}} - if {key: val, val: yo} map was provided then yo will replace {{key}}
// {#key}} hi {/key}} - hi will be rendered if key variable was provided as a template variable

func ReplaceVariable(input string, variableData map[string]interface{}) string {

	rootNode := NewBlockNode()
	parentNodes := []*BlockNode{rootNode}
	currentParentNode := rootNode

	codeReader := strings.NewReader(input)
	scn := bufio.NewScanner(codeReader)
	scn.Split(func (data []byte, atEOF bool) (advance int, token []byte, err error) {
		lenny := len(data)
		if atEOF && lenny == 0 {
			return 0, nil, nil
		}

		if data[0] == 123 {
			for i, val := range data {
				if val == 125 {
					if i +1 != lenny && data[i + 1] == 125 {
						return i + 2 , data[0:i+2], nil
					} else {
						return i + 1 , data[0:i+1], nil
					}
				}

				if val == 123 && i != 0 && i != 1 {
					return i  , data[0:i], nil
				}
			}
		} else {
			for i, val := range data {
				if val == 123 {
					return i  , data[0:i], nil
				}
			}
		}

		return 1, data[0:1], nil
	})

	var nodeType NodeType
	var offset int
	for scn.Scan() {
		byties := scn.Bytes()
		bytiesLen := len(byties)
		if byties[0] == 123 && byties[bytiesLen-1] == 125 {
			//fmt.Sprintf("%v", byties[1], "YYYYY")
			switch byties[1] {
			case Brace:
				nodeType = DoubleAction
				offset = 2

			case FSlash:
				nodeExpression := string(byties[2:bytiesLen-2])
				if nodeExpression != currentParentNode.Expression() {
					panic("Expected end block for " + currentParentNode.Expression() + " found end block for " + nodeExpression)
				}

				parentNodes = parentNodes[0:len(parentNodes)-1]
				currentParentNode = parentNodes[len(parentNodes)-1]

				continue

			case Hash:
				nodeType = BlockTextNode
				offset = 2

				node := NewBlockNode(WithDict(variableData),
									WithExpression((string(byties[offset:bytiesLen-offset]))))

				currentParentNode.AddChild(node)
				parentNodes = append(parentNodes, node)
				currentParentNode = node

				continue

			default:
				nodeType = SingleAction
				offset = 1
			}

			node := NewActionNode(WithActionType(nodeType),
								WithActionDictionary(variableData),
								WithActionText((string(byties[offset:bytiesLen-offset]))))

			currentParentNode.AddChild(node)

		} else {
			node := NewTextNode(WithText(string(byties)))
			currentParentNode.AddChild(node)
		}

	}

	//for _, val := range nodes {
	//	fmt.Println("|" + val.String() + "|")
	//}

	return rootNode.String()
}
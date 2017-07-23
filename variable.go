package variable

import (
	"strings"
	"text/scanner"
	"fmt"
	"net/http"

	"encoding/json"
	"bytes"
	"errors"
	"io/ioutil"
)

// TODO: implement replace variables functions.
// Variables are valid json, they can be either an array or object.


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

func getValueFromMap(mappy map[string]interface{}, key string) (string, error) {
	keys := strings.Split(key, ".")

	return getValueFromMap1(mappy, keys)
}

func getValueFromMap1(mappy map[string]interface{}, key []string) (string, error) {
	if len(key) > 1 {
		if _, ok := mappy[key[0]]; !ok {
			return ``, errors.New("value for key not found")
		}

		return getValueFromMap1(mappy, key[1:])
 	} else {
		if val, ok := mappy[key[0]]; ok {
			return fmt.Sprintf("%v", val), nil
		} else {
			return ``, errors.New("value for key not found")
		}
	}
}



func ReplaceVariable(input string, variableData map[string]interface{}) string {
	var scn scanner.Scanner
	var ret bytes.Buffer
	var nodes []Node



	codeReader := strings.NewReader(input)

	scn.Init(codeReader)
	tok := scn.Scan()

	inAction := false
	action := ``
	nodeType := NormalNode
	for tok != scanner.EOF {
		if tok == '{' {
			switch scn.Peek() {
			case '{':
				nodeType = DoubleAction
				scn.Scan()

			case '#':
				scn.Scan()
				nodeType = BlocAction


			case '/':


			default:
				nodeType = SingleAction
			}

			inAction = true

		} else if tok == '}' {
			if !inAction {
				node := NewTextNode(WithText(string(tok) + scn.TokenText()))

				nodes = append(nodes, node)

			} else {
				if nodeType == DoubleAction {
					scn.Scan()
				}
				node := NewActionNode(WithActionType(nodeType),
										WithActionText(action),
										WithActionDictionary(variableData))

				nodes = append(nodes, node)
			}

			inAction = false
			action = ``

		} else if inAction {
			action += scn.TokenText()

		} else {
			node := NewTextNode(WithText(scn.TokenText()))
			nodes = append(nodes, node)
		}

		tok = scn.Scan()
	}
	//
	//for _, val := range nodes {
	//	fmt.Println(val)
	//}

	for _, node := range nodes {
		ret.WriteString(node.String())
	}


//	fmt.Println("RERERERER", input, ret.String())

	return ret.String()

}
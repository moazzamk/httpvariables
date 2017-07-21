package variable

import (
	"strings"
	"text/scanner"
	"fmt"
	"net/http"

	"encoding/json"
	"github.com/pytlesk4/m"
	"bytes"
	"errors"
	"io/ioutil"
)

// TODO: implement replace variables functions.
// Variables are valid json, they can be either an array or object.


func PopulateRequestTemplate(req *http.Request, variables string) {
	if variables == `` {
		fmt.Println("EMPTTYTYTY")
		return
	}

	variableData := make(map[string]interface{})
	err := json.Unmarshal([]byte(variables), &variableData)
	if err != nil {
		fmt.Println("JSON MARSHAL ERROR ", err)
		fmt.Println("EMPTTYTYTY2")
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

	codeReader := strings.NewReader(input)

	scn.Init(codeReader)
	tok := scn.Scan()

	inAction := false
	action := ``
	for tok != scanner.EOF {

		if tok == '{' {
			inAction = true
			tok = scn.Scan()
			continue

		} else if (tok == '}') {
			if inAction {
				inAction = false

				val := m.Get(variableData, action)
				val1, err := getValueFromMap(variableData, action)
				fmt.Println("HHHHH", val, `|`, val1, `|`, err, `|`, action, `|VD: `, variableData)

				if val == nil {
					ret.WriteString("{" + action + "}")
				} else {
					ret.WriteString(fmt.Sprintf("%v", val))
				}
			}

			action = ``
			inAction = false

		} else if inAction {
			action += scn.TokenText()

		} else {
			ret.WriteString(scn.TokenText())
		}

		tok = scn.Scan()
	}



	fmt.Println("RERERERER", ret.String())

	return ret.String()

}
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

	req.Method = replaceVariable(req.Method, variableData)
	req.URL.Host = replaceVariable(req.URL.Host, variableData)
	req.URL.RawQuery = replaceVariable(req.URL.RawQuery, variableData)
	req.URL.Path = replaceVariable(req.URL.Path, variableData)
}

func getValueFromMap(mappy map[string]interface{}, key string) (string, error) {
	keys := strings.Split(key, ".")

	return getValueFromMap1(mappy, keys)
}

func getValueFromMap1(mappy map[string]interface{}, key []string) (string, error) {
	if len(key) > 1 {
		if _, ok := mappy[key[0]]; !ok {
			return nil, errors.New("value for key not found")
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



func replaceVariable(input string, variableData map[string]interface{}) string {
	var scn scanner.Scanner
	var ret bytes.Buffer

	codeReader := strings.NewReader(input)

	scn.Init(codeReader)
	tok := scn.Scan()

	inAction := false
	for tok != scanner.EOF {

		if tok == '{' {
			inAction = true
			tok = scn.Scan()
			continue
		} else if (tok == '}') {
			inAction = false
			tok = scn.Scan()
			continue
		}

		if inAction {
			val := m.Get(variableData, scn.TokenText())
			if val == nil {
				fmt.Println("NOBODY", scn.TokenText(), variableData)
				ret.WriteString("{" + scn.TokenText() + "}")
			} else {
				fmt.Println("NOBODY2", val, scn.TokenText(), variableData)
				ret.WriteString(fmt.Sprintf("%v", val))
			}

		} else {
			ret.WriteString(scn.TokenText())
		}

		tok = scn.Scan()
	}

	return ret.String()

}
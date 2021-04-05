package util

import (
	"bytes"
	"fmt"
	"text/template"
)

func templatePrintfEx(format string, m map[string]interface{}) (string, error) {
	// :TODO: try to use text/template/parse instead
	tpl, err := template.New("").Parse(format)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, m)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}

// TemplatePrintf() is like fmt.Sprintf() but uses
// named variables with syntax {{.name}}; implemented via text/template engine
func TemplatePrintf(template string, m map[string]interface{}) string {
	s, err := templatePrintfEx(template, m)
	if err != nil {
		// like fmt.Printf, do not panic in favor of encapsulating error in the result
		s = fmt.Sprintf("template error: %s:\n%s:\n%v", err, template, m)
	}
	return s
}

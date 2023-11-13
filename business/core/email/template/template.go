// Package template render body email.
package template

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	errorssys "myc-devices-simulator/business/sys/errors"
)

//go:embed *
var Template embed.FS

const longLine = 127

// Render render template contend.
func Render(language, templateFile string, data interface{}) (*bytes.Buffer, error) {
	templTxt, err := getFileTemplateContext(language + "/" + templateFile)
	if err != nil {
		return nil, fmt.Errorf("core.email.template.Render.getFileTemplateContext(-) - error: {%w}", err)
	}

	if !validateTemplate(templTxt) {
		return nil, fmt.Errorf("core.email.template.Render.validateTemplate(-) "+
			"- error: invalid Template - mycError: {%w}", errorssys.ErrEmailRenderTemplate)
	}

	newTmpl := template.New(templateFile)

	newTmpl, err = newTmpl.Parse(templTxt)
	if err != nil {
		return nil, fmt.Errorf("template.Render.Parse(-) - error: {%w} mycError: {%w}", err, errorssys.ErrEmailRenderTemplate)
	}

	var body bytes.Buffer

	if err = newTmpl.Execute(&body, data); err != nil {
		return nil, fmt.Errorf("core.email.template.Render.Execute(-) "+
			"- error: {%v} - mycError: {%w}", err, errorssys.ErrEmailRenderTemplate)
	}

	return &body, nil
}

// getFileTemplateContext get context file of the template.
func getFileTemplateContext(fileName string) (string, error) {
	file, err := Template.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("core.email.template.getFileTemplateContext.ReadFile(%s) "+
			"- error: {%v} - mycError {%w}", fileName, err, errorssys.ErrEmailReadFileTemplate)
	}

	return string(file), nil
}

// validateTemplate validation template length of line.
func validateTemplate(tmplText string) bool {
	for _, r := range tmplText {
		if r > longLine {
			return false
		}
	}

	return true
}

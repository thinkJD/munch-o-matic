package core

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

// SendTemplateNotification takes a template and matching sate to render the template and
// send to the provided topic.
func SendTemplateNotification(Topic string, Template string, Data interface{}) error {
	tmpl, err := template.New("template").Parse(Template)
	if err != nil {
		return fmt.Errorf("create template: %w", err)
	}
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, Data)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	renderedTemplate := buf.String()
	err = SendNotification(Topic, renderedTemplate)
	if err != nil {
		return fmt.Errorf("sending failed: %w", err)
	}
	return nil
}

// SendNotification sends the provided message to the topic
func SendNotification(Topic string, Content string) error {
	http.Post("https://ntfy.sh/"+Topic, "text/markdown",
		strings.NewReader(Content))

	return nil
}

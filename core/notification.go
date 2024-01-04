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
func SendTemplateNotification(Topic string, Title string, Template string, Data interface{}) error {
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
	err = SendNotification(Topic, Title, renderedTemplate)
	if err != nil {
		return fmt.Errorf("sending failed: %w", err)
	}
	return nil
}

// SendNotification sends the provided message to the topic
func SendNotification(Topic string, Title string, Message string) error {
	req, _ := http.NewRequest("POST", "https://ntfy.sh/"+Topic, strings.NewReader(Message))
	req.Header.Set("Title", Title)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send notification: %w", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("send notification failed with status %s", resp.Status)
	}

	return nil
}

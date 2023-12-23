package core

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

func SendAccountBalanceNotification(AccountBalance int) error {
	tplString := `Account balance low: {{.AccountBalance}}`

	tmpl, err := template.New("accountBalance").Parse(tplString)
	if err != nil {
		return fmt.Errorf("create template: %w", err)
	}

	var buf bytes.Buffer

	data := map[string]int{
		"AccountBalance": AccountBalance,
	}

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	resultString := buf.String()

	err = SendNotification("thinkjd_munch_o_matic", resultString)
	if err != nil {
		return fmt.Errorf("Sending failed: %w", err)
	}
	return nil
}

func SendNotification(Topic string, Content string) error {
	http.Post("https://ntfy.sh/"+Topic, "text/markdown",
		strings.NewReader(Content))

	return nil
}

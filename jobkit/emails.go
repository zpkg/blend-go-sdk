package jobkit

import (
	"context"
	"time"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/template"
)

// SendEmail sends a job status email.
func SendEmail(ctx context.Context, sender email.Sender, jobStatus JobStatus, ji *cron.JobInvocation, options ...email.MessageOption) error {
	if sender == nil {
		return nil
	}
	message, err := NewEmailMessage(ji.Name, jobStatus, ji.Err, ji.Elapsed, options...)
	if err != nil {
		return err
	}
	return sender.Send(ctx, message)
}

// NewEmailMessage returns a new email message.
func NewEmailMessage(jobName string, jobStatus JobStatus, jobErr error, elapsed time.Duration, options ...email.MessageOption) (email.Message, error) {
	message := email.Message{}

	vars := map[string]interface{}{
		"jobName":   jobName,
		"jobStatus": string(jobStatus),
		"elapsed":   elapsed,
		"err":       jobErr,
	}
	var err error
	message.Subject, err = template.New().WithBody(DefaultEmailSubjectTemplate).WithVars(vars).ProcessString()
	if err != nil {
		return message, err
	}
	message.HTMLBody, err = template.New().WithBody(DefaultEmailHTMLBodyTemplate).WithVars(vars).ProcessString()
	if err != nil {
		return message, err
	}
	message.TextBody, err = template.New().WithBody(DefaultEmailTextBodyTemplate).WithVars(vars).ProcessString()
	if err != nil {
		return message, err
	}

	return email.ApplyMessageOptions(message, options...), nil
}

const (
	// DefaultEmailMimeType is the default email mime type.
	DefaultEmailMimeType = "text/plain"

	// DefaultEmailSubjectTemplate is the default subject template.
	DefaultEmailSubjectTemplate = `({{ .Var "environment" }}) {{.Var "jobName" }} :: {{ .Var "jobStatus" }}`

	// DefaultEmailHTMLBodyTemplate is the default email html body template.
	DefaultEmailHTMLBodyTemplate = `
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<title>{{ .Var "jobName" }} {{ .Var "jobStatus" "unknown" }}</title>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
<meta http-equiv="X-UA-Compatible" content="IE=edge" />
<meta name="viewport" content="width=device-width, initial-scale=1.0 " />
<style>
.email-body {
	margin: 0;
	padding: 20px;
	font-family: sans-serif;
	font-size: 16pt;
}
</style>
</head>
<body class="email-body">
	<h2>{{ .Var "jobName" }} {{ .Var "jobStatus" "Unknown" }}</h2>
	<div class="email-details">
	{{ if .Var "err" }}
	<pre>{{ .Var "err" }}</pre>
	{{ end }}
	</div>
</body>
</html>
`

	// DefaultEmailTextBodyTemplate is the default body template.
	DefaultEmailTextBodyTemplate = `{{ .Var "jobName" }} {{ .Var "jobStatus" }}
Elapsed: {{ .Var "result.elapsed" }}
{{ if .HasVar "err" }}Error: {{ .Var "err" }}{{end}}`
)

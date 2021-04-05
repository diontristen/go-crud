package util

import (
	"bytes"
	"net/smtp"
	"path/filepath"
	"strings"
	"text/template"
)

// SendEmail takes a html template placed inside the templates folder, replaces the variables on templateStruct interface and
// send an email to the users defined on to array
func SendEmail(to []string, templateName string, subject string, templateStruct interface{}, ac *AppContext) (err error) {
	cfg := ac.Config

	// TODO: move sensitive to a external config, context
	from := cfg.EmailFrom
	password := cfg.EmailPassword
	smtpHost := cfg.EmailSMTPHost
	smtpPort := cfg.EmailStmpPort
	userName := cfg.EmailSMTPUsername

	templatePath, templateError := filepath.Abs("./templates/" + templateName)

	if templateError != nil {
		return templateError
	}

	t, parseError := template.ParseFiles(templatePath)

	if parseError != nil {
		return parseError
	}

	var body bytes.Buffer

	t.Execute(&body, templateStruct)

	auth := smtp.PlainAuth("", userName, password, smtpHost)

	toHeader := strings.Join(to, ",")

	msg := []byte("To: " + toHeader + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"\r\n" +
		body.String() + "\r\n")

	smtpError := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if smtpError != nil {
		return smtpError
	}

	return nil

}

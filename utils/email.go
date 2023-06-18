package utils

import (
	"bytes"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

type EmailInfo struct {
	Name    string
	Details string
	Otp     string
	Subject string
}

func SendEmail(emailAddr string, username string, info EmailInfo, path string) (string, error) {

	configs, err := LoadConfig("./../")
	if err != nil {
		return "", err
	}

	// load from file path
	filePath := fmt.Sprintf("%s/html/template.html", path)

	tp, err := template.ParseFiles(filePath)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = tp.Execute(&tpl, struct {
		Name string
		Otp  string
	}{Name: info.Name, Otp: info.Otp})
	if err != nil {
		return "", err
	}

	fromEmail := configs.Email
	password := configs.EmailPass

	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", emailAddr)
	m.SetHeader("Subject", "OTP for DefiRaise")
	m.SetBody("text/html", tpl.String()) // attach whatever you want

	d := gomail.NewDialer("smtp.gmail.com", 465, fromEmail, password)

	if err := d.DialAndSend(m); err != nil {
		return "", err
	}

	return "Email Sent!", nil
}

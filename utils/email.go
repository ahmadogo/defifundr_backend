package utils

import (
	"bytes"
	"html/template"

	"gopkg.in/gomail.v2"
)

type info struct {
	Name    string
	Details string
}

func SendEmail(emailAddr string, username string, details string) (string, error) {

	configs, err := LoadConfig("./../")
	if err != nil {
		return "", err
	}

	// load from file path
	path := "./../files/html/template.html"


	tp, err := template.ParseFiles(path)

	if err != nil {
		return "", err
	}

	var i info

	i.Name = username
	i.Details = details

	var tpl bytes.Buffer
	err = tp.Execute(&tpl, struct {
		Name string
		Otp  string
	}{Name: "demola234", Otp: "123456"})
	if err != nil {
		return "", err
	}

	fromEmail := configs.Email
	password := configs.EmailPass

	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", emailAddr)
	m.SetHeader("Subject", "golang test")
	m.SetBody("text/html", tpl.String()) // attach whatever you want

	d := gomail.NewDialer("smtp.gmail.com", 465, fromEmail, password)

	if err := d.DialAndSend(m); err != nil {
		return "", err
	}

	return "Email Sent!", nil
}

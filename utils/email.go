package utils

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
)

func SendEmail() {
	e := &email.Email{
		To:      []string{"test@example.com"},
		From:    "Jordan Wright <test@gmail.com>",
		Subject: "Awesome Subject",
		Text:    []byte("Text Body is, of course, supported!"),
		HTML:    []byte("<h1>Fancy HTML is supported, too!</h1>"),
		Headers: textproto.MIMEHeader{},
	}
	e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "junwu.shao@han-networks.com", "", "smtp.gmail.com"))
}

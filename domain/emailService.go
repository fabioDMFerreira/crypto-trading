package domain

import "net/smtp"

type SendMail func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

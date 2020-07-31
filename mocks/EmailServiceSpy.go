package mocks

import "net/smtp"

type EmailServiceSpy struct {
	SendMailCalls [][]interface{}
}

func (e *EmailServiceSpy) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	e.SendMailCalls = append(e.SendMailCalls, []interface{}{addr, a, from, to, msg})
	return nil
}

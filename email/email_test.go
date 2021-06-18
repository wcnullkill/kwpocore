package email

import "testing"

var sr EmailService

func init() {
	sr = DefaultEmailService(host, user, pwd, port)
}

func TestDial(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fail()
		}
	}()
	sr = DefaultEmailService(host, user, pwd, port)
}

func TestEmailSend(t *testing.T) {
	msg := NewEmailMsg()
	msg.To([]string{"wangchi@ctrchina.cn"})
	msg.CC([]string{"zhaoran@ctrchina.cn"})
	msg.From("kwp.cn@ctrchina.cn")
	msg.Subject("test")
	msg.Body([]byte("test email"))
	msg.AttachFile("test.txt", []byte("test file"))
	err := sr.Send(msg)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

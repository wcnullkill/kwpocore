package email

import (
	"crypto/tls"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	Send(msg *EmailMsg) error
}

type goEmailService struct {
	d *gomail.Dialer
}

func DefaultEmailService(host, user, pwd string, port int) EmailService {
	d := gomail.NewDialer(host, port, user, pwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &goEmailService{
		d: d,
	}
}

func (sr *goEmailService) Send(msg *EmailMsg) error {

	m := gomail.NewMessage()
	if len(msg.to) > 0 {
		m.SetHeader("To", msg.to...)
	}
	if len(msg.subject) > 0 {
		m.SetHeader("Subject", msg.subject)
	}
	t := "text/html"
	c := ""
	if len(msg.contenttype) > 0 {
		t = msg.contenttype
	}
	if len(msg.body) > 0 {
		c = string(msg.body)
	}
	m.SetBody(t, c)

	if len(msg.attachments) > 0 {
		rand.Seed(time.Now().Unix())
		r := rand.Int()
		td, err := ioutil.TempDir(os.TempDir(), strconv.Itoa(r))
		if err != nil {
			return err
		}
		defer os.RemoveAll(td)

		for name, file := range msg.attachments {
			fullname := path.Join(td, name)
			err = ioutil.WriteFile(fullname, file, 0755)
			if err != nil {
				return err
			}
			m.Attach(fullname)
		}
	}
	from := "kwpo.cn@ctrchina.cn"
	if len(msg.from) > 0 {
		from = msg.from
	}
	m.SetHeader("From", from)
	m.SetHeader("Cc", msg.cc...)
	return sr.d.DialAndSend(m)
}

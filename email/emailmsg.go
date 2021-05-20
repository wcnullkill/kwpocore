package email

type EmailMsg struct {
	to          []string
	cc          []string
	from        string
	subject     string
	contenttype string
	body        []byte
	attachments map[string][]byte
}

func NewEmailMsg() *EmailMsg {
	return &EmailMsg{attachments: make(map[string][]byte)}
}
func (msg *EmailMsg) To(to []string) {
	msg.to = to
}
func (msg *EmailMsg) CC(cc []string) {
	msg.cc = cc
}
func (msg *EmailMsg) From(from string) {
	msg.from = from
}

func (msg *EmailMsg) Subject(sub string) {
	msg.subject = sub
}

func (msg *EmailMsg) ContentType(t string) {
	msg.contenttype = t
}
func (msg *EmailMsg) Body(b []byte) {
	msg.body = b
}

func (msg *EmailMsg) AttachFile(filename string, file []byte) {
	msg.attachments[filename] = file
}

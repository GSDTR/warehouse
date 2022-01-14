package email

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"io"
	"io/ioutil"
	"log"
)

func Read_messages(login string, password string, msgCnt uint32) []map[string]string {
	log.Println("Connecting to server...")

	res := make([]map[string]string, 1, 1)

	// Connect to server
	c, err := client.DialTLS("imap.mail.ru:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(login, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > msgCnt {
		// We're using unsigned integers here, only substract if the result is > 0
		from = mbox.Messages - msgCnt
	}

	seqSet := new(imap.SeqSet)
	//	seqSet.AddNum(mbox.Messages)
	seqSet.AddRange(from, to)

	// Get the whole message body
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	//msg := <-messages?
	for msg := range messages {
		cur := make(map[string]string)
		if msg == nil {
			log.Fatal("Server didn't returned message")
		}

		r := msg.GetBody(section)
		if r == nil {
			log.Fatal("Server didn't returned message body")
		}

		// Create a new mail reader
		mr, err := mail.CreateReader(r)
		if err != nil {
			log.Fatal(err)
		}

		// Print some info about the message
		header := mr.Header
		if date, err := header.Date(); err == nil {
			cur["date"] = date.String()
		}
		if from, err := header.AddressList("From"); err == nil {
			if len(from) != 0 {
				cur["from"] = from[0].String()
			}
		}
		if to, err := header.AddressList("To"); err == nil {
			if len(to) != 0 {
				cur["to"] = to[0].String()
			}
		}
		if subject, err := header.Subject(); err == nil {
			cur["subject"] = subject
		}

		// Process each message's part
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			switch h := p.Header.(type) {
			case mail.TextHeader:
				// This is the message's text (can be plain-text or HTML)
				b, _ := ioutil.ReadAll(p.Body)
				cur["body"] = string(b)
			case mail.AttachmentHeader:
				// This is an attachment
				filename, _ := h.Filename()
				//fmt.Println("\tGot attachment: ", filename)
				cur["attachment"] = filename
			}
		}
		//fmt.Println("msg: ", cur["subject"])
		res = append(res, cur)
	}
	return res
}


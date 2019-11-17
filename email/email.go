package email

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/grokify/html-strip-tags-go"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func GetNewMessages() []string {
	c := connecting()

	messages, section := selectNewMessage(c, os.Getenv("mail_box"))
	var messagesString []string

	for msg := range messages {
		messagesString = append(messagesString, getFormattedMsg(msg, section))
	}
	log.Print(messagesString)
	return messagesString
}

func connecting() *client.Client {
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS(os.Getenv("imap"), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")
	// Login
	if err := c.Login(os.Getenv("email"), os.Getenv("email_pass")); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	return c
}

func selectNewMessage(c *client.Client, mailBox string) (chan *imap.Message, imap.BodySectionName) {
	// Select mailbox
	mbox, err := c.Select(mailBox, false)
	if err != nil {
		log.Fatal(err)
	}
	if mbox.UnseenSeqNum == 0 {
		log.Fatal("No new message in mailbox")
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(mbox.UnseenSeqNum)

	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 100)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	if messages == nil {
		log.Fatal("Server didn't returned message")
	}

	return messages, section
}

func getFormattedMsg(msg *imap.Message, section imap.BodySectionName) string {
	if msg == nil {
		log.Fatal("Server didn't returned message")
	}
	r := msg.GetBody(&section)
	if r == nil {
		log.Fatal("Server didn't returned message body")
	}

	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Fatal(err)
	}
	var texts []string

	// Save info about the message
	header := mr.Header
	if date, err := header.Date(); err == nil {
		texts = append(texts, "Date: "+date.String())
	}
	if from, err := header.AddressList("From"); err == nil {
		texts = append(texts, "From: "+addressListToString(from))
	}
	if subject, err := header.Subject(); err == nil {
		texts = append(texts, "Subject: "+subject)
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
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, _ := ioutil.ReadAll(p.Body)
			texts = append(texts, "Body: "+strip.StripTags(string(b)))

		case *mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			texts = append(texts, "Files: "+filename)
		}
	}

	return strings.Join(texts, "\n")
}

func addressListToString(l []*mail.Address) string {
	formatted := make([]string, len(l))
	for i, a := range l {
		formatted[i] = a.String()
	}
	return strings.Join(formatted, ", ")
}

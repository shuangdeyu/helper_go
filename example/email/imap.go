package email

import (
	"bytes"
	"fmt"
	"github.com/emersion/go-imap"
	id "github.com/emersion/go-imap-id"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	mxkClient "github.com/mxk/go-imap/imap"
	"io"
	"io/ioutil"
	"log"
	net_mail "net/mail"
	"os"
	"time"
)

/**
 * 连接 imap 邮件服务器
 */
func ConnectImap() {
	server := "imap.163.com:143"
	//server := "imap.163.com:993"
	// 你的邮件账号和密码
	username := "xiaof_for@163.com"
	password := "PAQWUJKIFHKPLGPP"

	c, err := client.Dial(server)
	//c, err := client.DialTLS(server, nil)
	if err != nil {
		log.Println("client connect error: ", err)
	}
	//c.StartTLS()

	// 登陆
	if err = c.Login(username, password); err != nil {
		log.Println("login error: ", err)
	}

	idClient := id.NewClient(c)
	idClient.ID(
		id.ID{
			id.FieldName:    "IMAPClient",
			id.FieldVersion: "3.1.0",
		},
	)
	defer c.Close()

	// 邮箱文件夹列表
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("邮箱文件夹:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// 选择收件箱
	mbox, err := c.Select("INBOX", false)
	if err != nil {

		fmt.Println("select inbox err: ", err)
		return
	}
	fmt.Println(mbox)
	if mbox.Messages == 0 {
		return
	}

	fmt.Println("未读邮件数:", mbox.Recent)
	return

	// 获得最新的十封邮件
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 10 {
		from = mbox.Messages - 10
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	section := imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, items, messages)
	}()
	log.Println("最后十封邮件:")
	imap.CharsetReader = charset.Reader
	for msg := range messages {
		r := msg.GetBody(&section)
		if r == nil {
			log.Fatal("服务器未返回邮件正文")
		}
		mr, err := mail.CreateReader(r)
		if err != nil {
			log.Fatal(err)
		}

		header := mr.Header
		var subject string
		if date, err := header.Date(); err == nil {
			log.Println("Date:", date)
		}
		if from, err := header.AddressList("From"); err == nil {
			log.Println("From:", from)
		}
		if to, err := header.AddressList("To"); err == nil {
			log.Println("To:", to)
		}
		if subject, err = header.Subject(); err == nil {
			log.Println("Subject:", subject)
		}

		// 处理邮件正文
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal("NextPart:err ", err)
			}

			switch h := p.Header.(type) {
			case *mail.InlineHeader:
				// 正文消息文本
				b, _ := ioutil.ReadAll(p.Body)
				mailFile := fmt.Sprintf("INBOX/%s.eml", subject)
				f, _ := os.OpenFile(mailFile, os.O_RDWR|os.O_CREATE, 0766)
				f.Write(b)
				f.Close()
			case *mail.AttachmentHeader:
				// 正文内附件
				filename, _ := h.Filename()
				log.Printf("attachment: %v\n", filename)
			}
		}
	}

	// 选择收取邮件的时间段
	criteria := imap.NewSearchCriteria()
	// 收取7天之内的邮件
	t1, err := time.Parse("2006-01-02 15:04:05", "2020-03-02 15:04:05")
	criteria.Since = t1
	// 按条件查询邮件
	ids, err := c.Search(criteria)
	fmt.Println(ids)
	if err != nil {
		fmt.Println(err)
	}
	if len(ids) == 0 {
		return
	}

	seqsetm := new(imap.SeqSet)
	seqsetm.AddNum(ids...)
	sect := &imap.BodySectionName{}
	messagesm := make(chan *imap.Message, 100)
	donem := make(chan error, 1)
	go func() {
		donem <- c.Fetch(seqsetm, []imap.FetchItem{sect.FetchItem()}, messagesm)
	}()
	for msg := range messagesm {
		r := msg.GetBody(sect)
		m, err := message.Read(r)
		if err != nil {
			fmt.Println(err)
			// return err
		}
		header := m.Header
		emailDate, _ := net_mail.ParseDate(header.Get("Date"))
		// 读取邮件内容
		fmt.Println(header, emailDate)
	}
}

func ConnectImapByMxkClient() {
	//server := "imap.163.com:143"
	server := "imap.163.com:993"
	// 你的邮件账号和密码
	username := "xiaof_for@163.com"
	password := "PAQWUJKIFHKPLGPP"

	var (
		c   *mxkClient.Client
		cmd *mxkClient.Command
		rsp *mxkClient.Response
	)

	// Connect to the server
	//c, _ = imap.Dial(server)
	c, _ = mxkClient.DialTLS(server, nil)

	// Remember to log out and close the connection when finished
	defer c.Logout(30 * time.Second)

	// Print server greeting (first response in the unilateral server data queue)
	fmt.Println("Server says hello:", c.Data[0].Info)
	c.Data = nil

	// Enable encryption, if supported by the server
	if c.Caps["STARTTLS"] {
		c.StartTLS(nil)
	}

	// Authenticate
	if c.State() == mxkClient.Login {
		c.Login(username, password)
	}

	// List all top-level mailboxes, wait for the command to finish
	cmd, _ = mxkClient.Wait(c.List("", "%"))

	// Print mailbox information
	fmt.Println("\nTop-level mailboxes:")
	for _, rsp = range cmd.Data {
		fmt.Println("|--", rsp.MailboxInfo())
	}

	// Check for new unilateral server data responses
	for _, rsp = range c.Data {
		fmt.Println("Server data:", rsp)
	}
	c.Data = nil

	// Open a mailbox (synchronous command - no need for imap.Wait)
	cm, err := c.Select("INBOX", true)
	fmt.Println(cm, err)
	fmt.Print("\nMailbox status:\n", c.Mailbox)

	// Fetch the headers of the 10 most recent messages
	set, _ := mxkClient.NewSeqSet("")
	if c.Mailbox.Messages >= 10 {
		set.AddRange(c.Mailbox.Messages-9, c.Mailbox.Messages)
	} else {
		set.Add("1:*")
	}
	cmd, _ = c.Fetch(set, "RFC822.HEADER")

	// Process responses while the command is running
	fmt.Println("\nMost recent messages:")
	for cmd.InProgress() {
		// Wait for the next response (no timeout)
		c.Recv(-1)

		// Process command data
		for _, rsp = range cmd.Data {
			header := mxkClient.AsBytes(rsp.MessageInfo().Attrs["RFC822.HEADER"])
			if msg, _ := net_mail.ReadMessage(bytes.NewReader(header)); msg != nil {
				fmt.Println("|--", msg.Header.Get("Subject"))
			}
		}
		cmd.Data = nil

		// Process unilateral server data
		for _, rsp = range c.Data {
			fmt.Println("Server data:", rsp)
		}
		c.Data = nil
	}

	// Check command completion status
	if rsp, err := cmd.Result(mxkClient.OK); err != nil {
		if err == mxkClient.ErrAborted {
			fmt.Println("Fetch command aborted")
		} else {
			fmt.Println("Fetch error:", rsp.Info)
		}
	}
}

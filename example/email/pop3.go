package email

import (
	"bufio"
	"fmt"
	"github.com/knadh/go-pop3"
	"io"
	"log"
)

/**
 * 连接 pop3 邮件服务器
 */
func ConnectPop3() {
	username := "xiaof_for@163.com"
	password := "PAQWUJKIFHKPLGPP"

	server := "pop3.163.com"
	options := pop3.Opt{
		Host: server,
		//Port: 110,
		Port:       995,
		TLSEnabled: true,
	}
	p := pop3.New(options)

	// Create a new connection
	c, err := p.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	// Authenticate
	if err := c.Auth(username, password); err != nil {
		log.Fatal(err)
	}

	// Print the total number of messages and their size.
	count, size, _ := c.Stat()
	fmt.Println("total messages=", count, "size=", size)

	// Pull the list of all message IDs and their sizes.
	msgs, _ := c.List(0)
	for _, m := range msgs {
		fmt.Println("id=", m.ID, "size=", m.Size)
	}

	// Pull all messages on the server. Message IDs go from 1 to N.
	for id := 1; id <= count; id++ {
		m, _ := c.Retr(id)

		if id != 9 {
			continue
		}
		// 使用 bufio.Reader 逐行读取
		bufReader := bufio.NewReader(m.Body)
		for {
			line, err := bufReader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("Error reading:", err)
				return
			}
			fmt.Print(line)
		}
		break

		fmt.Println(id, "=", m.Header.Get("subject"))

		// To read the multi-part e-mail bodies, see:
		// https://github.com/emersion/go-message/blob/master/example_test.go#L12
	}

}

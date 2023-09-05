package models

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func Email(email string) {
	mailTo := []string{
		email,
	}
	//邮件主题
	subject := "GPT程序异常退出"
	// 邮件正文
	body := "<html><body><p>" +
		"由于特殊原因GPT程序异常退出！" +
		"</p></body></html>"

	err := SendMail(mailTo, subject, body)
	if err != nil {
		log.Println(err)
		fmt.Println("send fail")
	}
	fmt.Println("发送成功！")
}

func SendMail(mailTo []string, subject string, body string) error {
	mailConn := map[string]string{
		"user": viper.GetString("smtp.user"),
		"pass": viper.GetString("smtp.pass"),
		"host": viper.GetString("smtp.host"),
		"port": viper.GetString("smtp.port"),
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(mailConn["user"], "晚风漪GPT")) //这种方式可以添加别名，即“XX官方”
	m.SetHeader("To", mailTo...)                                     //发送给多个用户
	m.SetHeader("Subject", subject)                                  //设置邮件主题
	m.SetBody("text/html", body)                                     //设置邮件正文
	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err
}

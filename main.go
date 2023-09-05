package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/eatmoreapple/openwechat"
	"github.com/spf13/viper"
	"main.go/models"
	"main.go/utils"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			// 发生错误时执行特定操作，发送邮件
			fmt.Println("发生错误，执行特定操作...")
			fmt.Println(r)
			models.Email(viper.GetString("smtp.user"))
		}
	}()
	utils.DispositionInit()

	models.MYSQL()

	models.InitializeTheRoster()

	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		//判断是否触发关键字
		for _, v := range models.ToRemove {
			if strings.Contains(msg.Content, v) {
				if !msg.IsSendBySelf() {
					user, err := msg.Sender()
					if err != nil {
						log.Println("获取消息的发送者失败！")
					}
					models.PrivateAndGroup(msg, user)
					break
				}
			}
		}

	}
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	err := bot.Block()
	if err != nil {
		fmt.Println(err)
	}
}

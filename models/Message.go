package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/spf13/viper"
)

var (
	initiate           int                               //选择的是什么名单
	BlackAndWhiteLists = make(map[string]interface{}, 0) //黑白名单群名字
	Whitelist          = 1                               //白名单 1
	Blacklist          = 2                               //黑名单 2
	ToRemove           []string                          //AI关键字触发
)

// 初始化名单
func InitializeTheRoster() {
	initiate = viper.GetInt("initiate")
	for i, v := range viper.GetStringSlice("BlackAndWhiteLists") {
		BlackAndWhiteLists[v] = i
	}
	ToRemove = viper.GetStringSlice("AIkeywords")
}

// 判断私聊还是群聊
func PrivateAndGroup(msg *openwechat.Message, user *openwechat.User) {
	//判断是否获取到信息
	if user != nil {
		switch {
		//判断消息是否由群组发送
		case msg.IsSendByGroup() && !msg.IsSendBySelf():
			DisposeGroupList(msg, user)

		// 判断消息是否由好友发送
		case msg.IsSendByFriend():
			DisposePrivate(msg, user)

		//判断消息是否由自己发送
		case msg.IsSendBySelf():
		}
	}
}

// 处理群聊消息是否在黑白名单中
func DisposeGroupList(msg *openwechat.Message, user *openwechat.User) {
	//过滤指定群不触发
	//判断启动的是白名单还是黑名单
	_, exists := BlackAndWhiteLists[user.NickName]
	if initiate == Whitelist {
		if exists {
			user, _ = msg.SenderInGroup()
			DisposeGroup(msg, user)
		}

	} else if initiate == Blacklist {
		if !exists {
			user, _ = msg.SenderInGroup()
			DisposeGroup(msg, user)
		}
	}

}

// 处理群聊信息
func DisposeGroup(msg *openwechat.Message, user *openwechat.User) {
	//获取发送者的信息
	Group, err := msg.SenderInGroup()
	if err != nil {
		log.Println("获取群聊发送者信息出错！")
	}
	//打印发送的信息
	log.Println("Group: " + user.NickName + ">>>>" + Group.NickName + ":" + msg.Content)
	Function(msg, user, 2)

}

// 处理私聊消息
func DisposePrivate(msg *openwechat.Message, user *openwechat.User) {
	if !msg.IsArticle() {
		//打印发送的信息
		log.Println("user>>>>" + user.NickName + ":" + msg.Content)
		Function(msg, user, 1)
	}

}

// 挑选用户使用哪一个功能 number 1---私聊   2---群聊 .......
func Function(msg *openwechat.Message, user *openwechat.User, number int) {
	switch {
	case strings.Contains(msg.Content, "天气"):
		Filtration(msg, user, number)
	case strings.Contains(msg.Content, "星期") && strings.Contains(msg.Content, "天"):
		Time := time.Now()
		if strings.Contains(msg.Content, "今天") {
			// 确定星期几
			weekday := ChantheEnglish(Time.Weekday())
			context := fmt.Sprintf("日期(今天)：" + Time.Format("2006-01-02\n"+"星期："+weekday))
			Send(msg, user, context, number)
		} else if strings.Contains(msg.Content, "明天") {
			Time = Time.Add(24 * time.Hour)
			// 确定星期几
			weekday := ChantheEnglish(Time.Weekday())
			context := fmt.Sprintf("日期(明天)：" + Time.Format("2006-01-02\n"+"星期："+weekday))
			Send(msg, user, context, number)
		} else if strings.Contains(msg.Content, "后天") {
			Time = Time.Add(2 * 24 * time.Hour)
			// 确定星期几
			weekday := ChantheEnglish(Time.Weekday())
			context := fmt.Sprintf("日期(后天)：" + Time.Format("2006-01-02\n"+"星期："+weekday))
			Send(msg, user, context, number)
		}

	default:
		content := msg.Content
		for _, str := range ToRemove {
			content = strings.Replace(content, str, "", -1)
		}
		Send(msg, user, GPTmsg(content), number)
	}
}

// 发送信息
func Send(msg *openwechat.Message, user *openwechat.User, context string, number int) {
	switch number {
	case 1:
		//私聊
		msg.ReplyText(context)
	case 2:
		//群聊
		context = fmt.Sprintf("@" + user.NickName + "\u2005" + "\n" + context)
		msg.ReplyText(context)
	}
}

func ChantheEnglish(weekday time.Weekday) string {
	// 将星期几转换为中文
	weekdayStr := ""
	switch weekday {
	case time.Sunday:
		weekdayStr = "星期日"
	case time.Monday:
		weekdayStr = "星期一"
	case time.Tuesday:
		weekdayStr = "星期二"
	case time.Wednesday:
		weekdayStr = "星期三"
	case time.Thursday:
		weekdayStr = "星期四"
	case time.Friday:
		weekdayStr = "星期五"
	case time.Saturday:
		weekdayStr = "星期六"
	}
	return weekdayStr
}

package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/eatmoreapple/openwechat"
	"github.com/spf13/viper"
)

type WeatherResponse struct {
	Status   string `json:"status"`   // 状态码
	Count    string `json:"count"`    // 数据数量
	Info     string `json:"info"`     // 状态信息
	Infocode string `json:"infocode"` // 信息码

	Forecasts []Forecast `json:"forecasts"` // 天气预报信息
}

type Forecast struct {
	City       string `json:"city"`       // 城市名称
	Adcode     string `json:"adcode"`     // 区域编码
	Province   string `json:"province"`   // 省份名称
	Reporttime string `json:"reporttime"` // 报告时间
	Casts      []Cast `json:"casts"`      // 天气预报
}

type Cast struct {
	Date           string `json:"date"`            // 日期
	Week           string `json:"week"`            // 星期几
	Dayweather     string `json:"dayweather"`      // 白天天气状况
	Nightweather   string `json:"nightweather"`    // 晚上天气状况
	Daytemp        string `json:"daytemp"`         // 白天温度
	Nighttemp      string `json:"nighttemp"`       // 晚上温度
	Daywind        string `json:"daywind"`         // 白天风向
	Nightwind      string `json:"nightwind"`       // 晚上风向
	Daypower       string `json:"daypower"`        // 白天风力
	Nightpower     string `json:"nightpower"`      // 晚上风力
	DaytempFloat   string `json:"daytemp_float"`   // 白天温度的浮点数表示
	NighttempFloat string `json:"nighttemp_float"` // 晚上温度的浮点数表示
}

type Data struct {
	WeatherResponse
	Prompt string
}

type Message struct {
	ToUserName   string `xml:"ToUserName"`   // 开发者微信号
	FromUserName string `xml:"FromUserName"` // 发送方帐号（一个OpenID）
	CreateTime   int64  `xml:"CreateTime"`   // 消息创建时间 （整型）
	MsgType      string `xml:"MsgType"`      // text
	Content      string `xml:"Content"`      // 文本消息内容
	MsgId        int64  `xml:"MsgId"`        // 消息id，64位整型
}

// 获取天气
func GetTheWeather(city string) Data {
	// 高德天气API接口地址
	apiUrl := "https://restapi.amap.com/v3/weather/weatherInfo"

	// 构造请求参数
	params := url.Values{}
	params.Set("key", viper.GetString("weather.key"))
	params.Set("city", city)
	params.Set("extensions", "all")

	// 发送HTTP GET请求
	resp, err := http.Get(apiUrl + "?" + params.Encode())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 解析JSON响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var weatherResp WeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		log.Fatal(err)
	}

	// 根据天气情况输出相应的提示
	var prompt string
	switch {
	case strings.Contains(weatherResp.Forecasts[0].Casts[0].Dayweather, "雨"):
		prompt = "可能会下雨，请注意带好雨具"
	case strings.Contains(weatherResp.Forecasts[0].Casts[0].Dayweather, "晴") && weatherResp.Forecasts[0].Casts[0].DaytempFloat >= "35":
		prompt = "天气晴朗，气温较高，请注意防晒和补水"

	case strings.Contains(weatherResp.Forecasts[0].Casts[0].Dayweather, "晴") && weatherResp.Forecasts[0].Casts[0].DaytempFloat <= "5":
		prompt = "天气晴朗，气温较低，请注意保暖"

	default:
		prompt = "天气不错，可以出门逛逛"

	}

	// 输出天气信息
	fmt.Printf("城市：%s\n", weatherResp.Forecasts[0].City)
	fmt.Printf("天气：%s\n", weatherResp.Forecasts[0].Casts[0].Dayweather)
	fmt.Printf("温度：%s℃\n", weatherResp.Forecasts[0].Casts[0].Daytemp)
	fmt.Printf("风力：%s级\n", weatherResp.Forecasts[0].Casts[0].Daypower)
	fmt.Printf("风向：%s\n", weatherResp.Forecasts[0].Casts[0].Daywind)
	fmt.Printf("发布时间：%s\n", weatherResp.Forecasts[0].Reporttime)
	fmt.Println(prompt)

	//将获取的信息传递过去
	data := Data{
		WeatherResponse: weatherResp,
		Prompt:          prompt,
	}

	return data
}

// 过滤不需要的字段只留城市字段
func Filtration(msg *openwechat.Message, user *openwechat.User, number int) {
	Content := msg.Content
	//获取过滤字段
	str := viper.GetStringSlice("Filtration")
	for _, c := range str {
		Content = strings.ReplaceAll(Content, c, "")
	}
	//查询数据库中的对应城市code编码
	usercode, addrs := FincCode(Content)
	if usercode != nil {
		data := GetTheWeather(usercode.Adcode)
		//根据用户发送的信息做出相应的处理
		UserRequirements(msg.Content, data, msg, number, user)
	} else {
		context := fmt.Sprintf("城市错误！请将“" + addrs + "”多余的修改后重新提交")
		Send(msg, user, context, number)

	}
}

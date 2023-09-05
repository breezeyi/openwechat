package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

type Response struct {
	Code    int      `json:"code"`
	Data    DataTime `json:"data"`
	Content string   `json:"content"`
	Sangbo  string   `json:"sangbo"`
}

type DataTime struct {
	Output string `json:"output"`
}

// GPT回复
func GPTmsg(content string) string {
	fmt.Println(">>>>>>>>>msg")
	apiUrl := "https://api.lolimi.cn/api/ai/a"

	// 构造请求参数
	params := url.Values{}
	params.Set("msg", content)
	params.Set("key", viper.GetString("GPTkey"))

	// 发送HTTP GET请求
	resp, err := http.Get(apiUrl + "?" + params.Encode())
	if err != nil {
		log.Println(err)
		return "未知错误！"
	}
	defer resp.Body.Close()

	// 解析JSON响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "未知错误！"
	}
	//fmt.Println(string(body))
	var msg Response
	err = json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Println(err)
		return "未知错误！"
	}
	return msg.Data.Output
}

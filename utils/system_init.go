package utils

import (
	"log"

	"github.com/spf13/viper"
)

// 初始化配置文件
func DispositionInit() {
	viper.SetConfigFile("./config/conf.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("读取配置文件错误！" + err.Error())
	} else {
		log.Println("配置文件加载成功！")
	}
}

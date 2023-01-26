package function

import (
	"fmt"
	"gopkg.in/ini.v1"
)

func readcfg(name string) string {
	var cfg *ini.File
	cfg, err := ini.Load("./cfg.ini")
	if err != nil {
		fmt.Println("read config error")
	}
	return cfg.Section("info").Key(name).String()
}

var (
	APPID            = readcfg("APPID")
	APPSECRET        = readcfg("APPSECRET")
	WeatTemplateID   = readcfg("WeatTemplateID") //天气模板ID，替换成自己的
	WeatherVersion   = readcfg("WeatherVersion")
	WeatherAppid     = readcfg("WeatherAppid")
	WeatherAppsecret = readcfg("WeatherAppsecret")
	Apikey           = readcfg("Apikey")
	rateTemplateID   = readcfg("rateTemplateID")
)

type token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

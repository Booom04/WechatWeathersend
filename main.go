package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"net/http"
	"strings"
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
	weatherAppid     = readcfg("weatherAppid")
	weatherAppsecret = readcfg("weatherAppsecret")
)

type token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// 发送天气预报
func weather() {
	access_token := getaccesstoken()
	fmt.Println(access_token)
	if access_token == "" {
		return
	}

	flist := getflist(access_token)
	if flist == nil {
		return
	}
	fmt.Println(flist)
	var city string
	for _, v := range flist {
		switch v.Str {
		case "**************":
			city = "**"
			go sendweather(access_token, city, v.Str)
		case "***************":
			city = "****"
			go sendweather(access_token, city, v.Str)
		case "****************":
			city = "***"
			go sendweather(access_token, city, v.Str)
		default:
			fmt.Println("没找到他们")
		}
	}
	fmt.Println("weather is ok")
}

// 获取微信access_token
func getaccesstoken() string {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", APPID, APPSECRET)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取微信token失败", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("微信token读取失败", err)
		return ""
	}

	token := token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println("微信token解析json失败", err)
		return ""
	}
	return token.AccessToken
}

// 获取关注者列表
func getflist(access_token string) []gjson.Result {
	url := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + access_token + "&next_openid="
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取关注列表失败", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取内容失败", err)
		return nil
	}
	flist := gjson.Get(string(body), "data.openid").Array()
	return flist
}

// 发送模板消息
func templatepost(access_token string, reqdata string, fxurl string, templateid string, openid string) {
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + access_token

	reqbody := "{\"touser\":\"" + openid + "\", \"template_id\":\"" + templateid + "\", \"url\":\"" + fxurl + "\", \"data\": " + reqdata + "}"

	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(string(reqbody)))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

// 获取天气
func getweather(city string) (string, string, string, string) {
	url := fmt.Sprintf("http://v0.yiketianqi.com/free/day?appid=%v&appsecret=%v&version=%v&city=%v", weatherAppid, weatherAppsecret, WeatherVersion, city)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取天气失败", err)
		return "", "", "", ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取内容失败", err)
		return "", "", "", ""
	}

	//data := gjson.Get(string(body), "data").Array()
	thisday := string(body)
	day := gjson.Get(thisday, "date").Str
	wea := gjson.Get(thisday, "wea").Str
	tem := gjson.Get(thisday, "tem").Str
	//tem2 := gjson.Get(thisday, "tem2").Str
	air_tips := gjson.Get(thisday, "air").Str
	return day, wea, tem, air_tips
}

// 发送天气
func sendweather(access_token, city, openid string) {
	day, wea, tem, air_tips := getweather(city)
	if day == "" || wea == "" || tem == "" || air_tips == "" {
		return
	}
	reqdata := "{\"city\":{\"value\":\"城市：" + city + "\", \"color\":\"#0000CD\"}, \"day\":{\"value\":\"" + day + "\"}, \"wea\":{\"value\":\"天气：" + wea + "\"}, \"tem1\":{\"value\":\"平均温度：" + tem + "\"}, \"air_tips\":{\"value\":\"tips：" + air_tips + "\"}}"
	//fmt.Println(reqdata)
	templatepost(access_token, reqdata, "", WeatTemplateID, openid)
}

func main() {
	//t := getaccesstoken()
	//sendweather(t, "石家庄", "o6KZd5jLDiiIskxiPXpj_X-mI8NI")
	weather()
	select {}

}

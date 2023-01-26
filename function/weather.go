package function

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

// 发送天气预报
func Weather() {
	access_token := getaccesstoken()
	//fmt.Println(access_token)
	if access_token == "" {
		return
	}

	flist := getflist(access_token)
	if flist == nil {
		return
	}
	//fmt.Println(flist)
	var city string
	for _, v := range flist {
		switch v.Str {
		case "*************************":  //填写openid
			city = "**"
			go sendweather(access_token, city, v.Str)
		//case "*****************":
		//	city = "**"
		//	go sendweather(access_token, city, v.Str)
		//case "********************":
		//	city = "**"
		//	go sendweather(access_token, city, v.Str)
		//case "**********************":
		//	city = "**"
		default:
		}
	}
	fmt.Println("weather is ok")
}

// 获取天气
func getweather(city string) (string, string, string, string) {
	url := fmt.Sprintf("http://v0.yiketianqi.com/free/day?appid=%v&appsecret=%v&version=%v&city=%v", WeatherAppid, WeatherAppsecret, WeatherVersion, city)
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

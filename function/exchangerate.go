package function

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

// 发送汇率信息
func Rate() {
	access_token := getaccesstoken()
	if access_token == "" {
		return
	}

	flist := getflist(access_token)
	if flist == nil {
		return
	}

	var currency string
	for _, v := range flist {
		switch v.Str {
		case "******************":
			currency = "USDT"
			go sendrate(access_token, currency, v.Str)
		default:

		}
	}
	fmt.Println("rate is ok")
}

// 获取汇率信息
func getrate(currency string) string {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=%v&to_currency=CNY&apikey=%v", currency, Apikey)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取汇率信息失败", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取内容失败", err)
		return ""
	}
	data := gjson.Get(string(body), "Realtime Currency Exchange Rate").Get("5\\. Exchange Rate").Str
	return data

}

// 发送汇率
func sendrate(access_token, currency, openid string) {
	Usdt := getrate(currency)
	if Usdt == "" {
		return
	}
	reqdata := "{\"Usdt\":{\"value\": " + Usdt + "}}"
	templatepost(access_token, reqdata, "", rateTemplateID, openid)
}

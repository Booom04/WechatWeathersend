package main

import (
	"WechatWeathersend/function"
	"fmt"
	"github.com/robfig/cron"
)

func main() {
	spec := "0 0 7 * * *"   //每天早晨7:00
	spec1 := "0 0 12 * * *" // 每天7:00
	spec2 := "0 0 18 * * *"
	c := cron.New()
	c.AddFunc(spec, function.Weather)
	c.AddFunc(spec, function.Rate)
	c.AddFunc(spec1, function.Rate)
	c.AddFunc(spec2, function.Rate)
	c.Start()
	fmt.Println("开启定时任务")
	select {}
	//function.Rate()
	//function.Weather()
}

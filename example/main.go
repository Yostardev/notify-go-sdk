package main

import "github.com/Yostardev/notify-go-sdk"

func SendEmail() error {
	return notify_go_sdk.NewClient("mock token").
		AddTemplateContentVariables("name", "1").
		AddTemplateContentVariables("xxx", "2").
		SetEnv(notify_go_sdk.Dev).
		SetTemplateId(10).
		AddReceiver("xx.xx@yo-star.com").
		Send()
}

func SendEmailByTemplateName() error {
	return notify_go_sdk.NewClient("mock token").
		SetTemplateContentVariables(map[string]string{
			"name": "1",
			"xxx":  "2",
		}).
		SetEnv(notify_go_sdk.Dev).
		SetTemplateName("测试HTML邮件模板").
		SetReceiver([]string{"xx.xx@yo-star.com"}).
		Send()
}

func SendEmailError() error {
	return notify_go_sdk.NewClient("mock token").
		SetEnv(notify_go_sdk.Dev).
		SetTemplateId(10).
		Send()
}

func main() {
	err := SendEmail()
	if err != nil {
		panic(err)
	}

	err = SendEmailByTemplateName()
	if err != nil {
		panic(err)
	}

	err = SendEmailError()
	if err != nil {
		panic(err)
	}
}

package notify_go_sdk

import (
	"errors"
	"github.com/Yostardev/gf"
	"github.com/Yostardev/requests"
	"strconv"
)

type Client struct {
	token       string
	env         EnvName
	deliverInfo DeliverInfo
}

type EnvName string

var Dev EnvName = "dev"
var Test EnvName = "test"
var Uat EnvName = "uat"
var Prod EnvName = "prod"

type DeliverInfo struct {
	TemplateId               int64             `json:"templateId"`               // 模板ID
	TemplateName             string            `json:"templateName"`             // 模板名称
	TemplateContentVariables map[string]string `json:"templateContentVariables"` // 模板内容变量
	Receiver                 []string          `json:"receiver"`                 // 接收者
}

type templateList struct {
	Code int `json:"code"`
	Data []struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

type deliverResponse struct {
	Code int `json:"code"`
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

func (c *Client) SetEnv(env EnvName) *Client {
	c.env = env
	return c
}

func (c *Client) SetTemplateId(templateId int64) *Client {
	c.deliverInfo.TemplateId = templateId
	return c
}

func (c *Client) SetTemplateName(templateName string) *Client {
	c.deliverInfo.TemplateName = templateName
	return c
}

func (c *Client) SetTemplateContentVariables(templateContentVariables map[string]string) *Client {
	c.deliverInfo.TemplateContentVariables = templateContentVariables
	return c
}

func (c *Client) SetReceiver(receiver []string) *Client {
	c.deliverInfo.Receiver = receiver
	return c
}

func (c *Client) AddTemplateContentVariables(key, value string) *Client {
	if c.deliverInfo.TemplateContentVariables == nil {
		c.deliverInfo.TemplateContentVariables = make(map[string]string)
	}
	c.deliverInfo.TemplateContentVariables[key] = value
	return c
}

func (c *Client) AddReceiver(receiver string) *Client {
	if c.deliverInfo.Receiver == nil {
		c.deliverInfo.Receiver = make([]string, 0)
	}
	c.deliverInfo.Receiver = append(c.deliverInfo.Receiver, receiver)
	return c
}

func (c *Client) getUri() string {
	switch c.env {
	case Dev:
		return "https://dev-notify-api.yostar.net"
	case Test:
		return "https://test-notify-api.yostar.net"
	case Uat:
		return "https://uat-notify-api.yostar.net"
	default:
		return "https://notify-api.yostar.net"
	}
}

func (c *Client) getTemplateIdByName(templateName string) (int64, error) {
	req := requests.New()
	req.AddHeader("Authorization", c.token)
	req.AddQuery("name", templateName)
	req.SetUrl(gf.StringJoin(c.getUri(), "/api/v1/template/list"))
	res, err := req.Get()
	if err != nil {
		return 0, err
	}
	if res.StatusCode != 200 {
		return 0, errors.New(gf.StringJoin("请求失败，状态码：", strconv.Itoa(res.StatusCode), "，错误信息：", res.Body.String()))
	}

	var templateList templateList
	err = res.Body.JsonBind(&templateList)
	if err != nil {
		return 0, err
	}

	if templateList.Code != 200 {
		return 0, errors.New(gf.StringJoin("请求失败，错误码：", strconv.Itoa(templateList.Code), "，错误信息：", res.Body.String()))
	}

	for i := range templateList.Data {
		if templateList.Data[i].Name == templateName {
			return templateList.Data[i].Id, nil
		}
	}
	return 0, errors.New("模板不存在")
}

func (c *Client) Send() error {
	if c.deliverInfo.TemplateId == 0 {
		templateId, err := c.getTemplateIdByName(c.deliverInfo.TemplateName)
		if err != nil {
			return err
		}
		c.deliverInfo.TemplateId = templateId
	}

	req := requests.New()
	req.AddHeader("Authorization", c.token)
	req.SetJsonBody(c.deliverInfo)
	req.SetUrl(gf.StringJoin(c.getUri(), "/api/v1/deliver"))
	res, err := req.Post()
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(gf.StringJoin("请求失败，状态码：", strconv.Itoa(res.StatusCode), "，错误信息：", res.Body.String()))
	}

	var deliverResponse deliverResponse
	err = res.Body.JsonBind(&deliverResponse)
	if err != nil {
		return err
	}

	if deliverResponse.Code != 200 {
		return errors.New(gf.StringJoin("请求失败，错误码：", strconv.Itoa(deliverResponse.Code), "，错误信息：", res.Body.String()))
	}
	return nil
}

package email

import (
	"backend-rms/core/config"
	"backend-rms/core/redis"
	"backend-rms/util/crypto"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const SendCLoudCommonURL = "https://api.sendcloud.net/apiv2/mail/send"
const BrevoCommonURL = "https://api.brevo.com/v3/smtp/email"

type EmailSender interface {
	SetParam(string, string, string, int) url.Values
	SendCommonEmail(url.Values, string) error
}

type SetParamReturn struct {
	SendCloudReturn url.Values
	BrevoReturn     BrevoRequest
}

type SendCloud struct {
}

type BrevoRequest struct {
	Sender      BrevoAddress   `json:"sender"`
	To          []BrevoAddress `json:"to"`
	Subject     string         `json:"subject"`
	HtmlContent string         `json:"htmlContent"`
}

type BrevoAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type Brevo struct {
}

type SendCloudReturn struct {
	Info       map[string][]string
	Message    string
	Result     bool
	StatusCode int
}

type BrevoReturn struct {
	MessageId string
	Message   string
	Code      string
}

func SendException(errString string) {
	noRepeat, _ := DoNotSendRepeat(errString)
	if !noRepeat {
		return
	}

	to := config.Conf.Email.DeveloperEmail
	title := config.Conf.Email.FromName + "-RMS系统异常"
	hostName, _ := os.Hostname()

	html := "Trigger Time：" + time.Now().String() + "<br/>" +
		"Host：" + hostName + "<br/>"

	html = html + "<br/>" + errString

	var setParamReturn SetParamReturn
	var urlAddress string
	var sendErr error
	if config.Conf.Email.EmailDriver == "sendCloud" {
		emailKind := new(SendCloud)
		setParamReturn = emailKind.SetParam(to, title, html, 0)
		urlAddress = SendCLoudCommonURL
		sendErr = emailKind.SendCommonEmail(setParamReturn.SendCloudReturn, urlAddress)
	} else {
		emailKind := new(Brevo)
		setParamReturn = emailKind.SetParam(to, title, html, 0)
		urlAddress = BrevoCommonURL
		sendErr = emailKind.SendCommonEmail(setParamReturn.BrevoReturn, urlAddress)
	}
	if sendErr != nil {
		panic("Send exception email is failed; email: " + to + ", error is: " + fmt.Sprintf("%v", sendErr))
	}
	return
}

func SendHttpException(err string, c *gin.Context) {
	noRepeat, _ := DoNotSendRepeat(err)
	if !noRepeat {
		return
	}
	to := config.Conf.Email.DeveloperEmail
	title := config.Conf.Email.FromName + "-RMS系统异常"
	hostName, _ := os.Hostname()
	urlName := c.Request.Host + c.Request.RequestURI

	html := "Trigger Time：" + time.Now().String() + "<br/>" +
		"Host：" + hostName + "<br/>" +
		"Url：" + urlName + "<br/>"

	var param string
	if c.Request.Method == "POST" {
		requestBody, _ := io.ReadAll(c.Request.Body)
		paramByte, _ := json.Marshal(requestBody)
		param = string(paramByte)
	}

	html = html + "Method：" + c.Request.Method + "<br/>" +
		"Param：" + param + "<br/>" +
		err

	var setParamReturn SetParamReturn
	var urlAddress string
	var sendErr error
	if config.Conf.Email.EmailDriver == "sendCloud" {
		emailKind := new(SendCloud)
		setParamReturn = emailKind.SetParam(to, title, html, 0)
		urlAddress = SendCLoudCommonURL
		sendErr = emailKind.SendCommonEmail(setParamReturn.SendCloudReturn, urlAddress)
	} else {
		emailKind := new(Brevo)
		setParamReturn = emailKind.SetParam(to, title, html, 0)
		urlAddress = BrevoCommonURL
		sendErr = emailKind.SendCommonEmail(setParamReturn.BrevoReturn, urlAddress)
	}
	if sendErr != nil {
		panic("Send exception email is failed; email: " + to + ", error is: " + fmt.Sprintf("%v", sendErr))
	}
	return
}

func (s *SendCloud) SendCommonEmail(urlValues url.Values, urlAddress string) error {
	res, err := http.PostForm(urlAddress, urlValues)
	if err != nil {
		return errors.Wrap(err, "send developer email failed")
	}
	defer res.Body.Close()

	BodyByte, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "io read all failed")
	}

	sendCloudReturn := SendCloudReturn{}
	err = json.Unmarshal(BodyByte, &sendCloudReturn)
	if err != nil {
		return errors.New("json: unmarshal failed")
	}
	if sendCloudReturn.Result != true {
		return errors.New("sendCloud error message：%v" + sendCloudReturn.Message)
	}
	return nil
}

func (b *Brevo) SendCommonEmail(brevoRequest BrevoRequest, urlAddress string) error {
	body, err := json.Marshal(brevoRequest)
	// 创建请求对象
	req, err := http.NewRequest("POST", urlAddress, bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "create request failed")
	}

	// 设置header
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", config.Conf.Email.ApiKey)

	// 发送请求
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "send developer email failed")
	}
	defer res.Body.Close()

	// 读取响应体
	BodyByte, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "io read all failed")
	}

	brevoReturn := BrevoReturn{}
	err = json.Unmarshal(BodyByte, &brevoReturn)
	if err != nil {
		return errors.New("json: unmarshal failed")
	}
	if brevoReturn.MessageId == "" {
		return errors.New("Brevo error message：" + brevoReturn.Message + "; code: " + brevoReturn.Code + ";")
	}
	return nil
}

func (s *SendCloud) SetParam(to string, title string, html string, emailType int) SetParamReturn {
	urlValues := url.Values{}
	urlValues.Add("apiKey", config.Conf.Email.ApiKey)
	//emailType:0触发邮件，1批量邮件
	if emailType == 0 {
		urlValues.Add("apiUser", config.Conf.Email.ApiUser)
		urlValues.Add("from", config.Conf.Email.From)
	} else {
		urlValues.Add("apiUser", config.Conf.Email.ApiUserBatch)
		urlValues.Add("from", config.Conf.Email.FromBatch)
	}
	urlValues.Add("to", to)
	urlValues.Add("subject", title)
	urlValues.Add("html", html)
	urlValues.Add("fromName", config.Conf.Email.FromName)
	return SetParamReturn{
		SendCloudReturn: urlValues,
		BrevoReturn:     BrevoRequest{},
	}
}

func (b *Brevo) SetParam(to string, title string, html string, emailType int) SetParamReturn {
	//emailType:0触发邮件，1批量邮件
	sender := BrevoAddress{}
	if emailType == 0 {
		sender.Email = config.Conf.Email.From
	} else {
		sender.Email = config.Conf.Email.FromBatch
	}
	sender.Name = config.Conf.Email.FromName

	toSlice := strings.Split(to, ";")
	var brevoTo []BrevoAddress
	for _, toAddress := range toSlice {
		if len(toAddress) != 0 {
			brevoTo = append(brevoTo, BrevoAddress{Email: toAddress})
		}
	}

	payload := BrevoRequest{
		Sender:      sender,
		To:          brevoTo,
		Subject:     title,
		HtmlContent: html,
	}
	return SetParamReturn{
		SendCloudReturn: nil,
		BrevoReturn:     payload,
	}
}

// 是否重复发送同一个邮件 同一个报错邮件每10分钟发送一次
func DoNotSendRepeat(errString string) (bool, error) {
	key := "error:" + crypto.Md5(errString)
	n, err := redis.RMSRedisCli.Exists(key)
	if err != nil {
		return false, errors.Wrap(err, "redis: get the key is exists failed")
	}
	if n == 1 {
		return false, nil
	}

	isSuccess, err := redis.RMSRedisCli.SetNX(key, true, 600*time.Second)
	if err != nil {
		return true, errors.Wrap(err, "redis: setNX key failed")
	}
	if isSuccess {
		return true, nil
	}
	return true, errors.New("setNX key failed")
}

package service

import (
	"backend-rms/core/config"
	"backend-rms/core/logger"
	"backend-rms/repo"
	"backend-rms/util"
	"backend-rms/util/email"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type emailService struct{}

var EmailService = &emailService{}

func (es *emailService) SendTemplate(event string, to string, param map[string]string, emailType int, worldCode string, userType int) {
	templates, err := getTemplates(event, param, worldCode, userType)
	if err != nil {
		logger.Errorf("get email template failed: %v\n", err)
	}

	for _, template := range templates {

		var setParamReturn email.SetParamReturn
		var urlAddress string
		if config.Conf.Email.EmailDriver == "sendCloud" {
			emailKind := new(email.SendCloud)
			setParamReturn = emailKind.SetParam(to, template["title"], template["html"], emailType)
			urlAddress = email.SendCLoudCommonURL
			err = emailKind.SendCommonEmail(setParamReturn.SendCloudReturn, urlAddress)
		} else {
			emailKind := new(email.Brevo)
			setParamReturn = emailKind.SetParam(to, template["title"], template["html"], emailType)
			urlAddress = email.BrevoCommonURL
			err = emailKind.SendCommonEmail(setParamReturn.BrevoReturn, urlAddress)
		}
		if err != nil {
			logger.Errorf("Send template email is failed; email: %v, event: %v, world code: %v，user type: %v, error is: %v", to, event, worldCode, userType, err)
			continue
		}
		logger.Infof("Send template email is success; email: %v, event: %v, world code: %v，user type: %v", to, event, worldCode, userType)
	}
	return
}

func getTemplates(event string, params map[string]string, worldCode string, userType int) ([]map[string]string, error) {
	//userType：1trial，2live
	//worldCode：默认发送国家是巴哈马
	if worldCode == "" {
		worldCode = "BS"
	}
	templateInfos, err := repo.EmailTemplateRepo.GetInfosByEvent(event, userType, worldCode)
	if err != nil {
		return nil, err
	}
	if len(templateInfos) == 0 {
		return nil, errors.New("not found " + event + "email template")
	}

	templates := make([]map[string]string, 0)
	//邮件模版先去查找only_world是否有合适的, 没有发送默认的, 即only_world为空的模版
	for _, templateInfo := range templateInfos {
		sliceOnlyWorld := strings.Split(templateInfo.OnlyWorld, ",")
		if len(sliceOnlyWorld) > 0 && util.InSlice(worldCode, sliceOnlyWorld) {
			html, err := replaceParams(params, templateInfo.Content, templateInfo.FrameId)
			if err != nil {
				return nil, err
			}
			templates = append(templates, map[string]string{"title": templateInfo.Title, "html": html})
		}
	}

	if len(templates) < 1 {
		for _, templateInfo := range templateInfos {
			sliceOnlyWorld := strings.Split(templateInfo.OnlyWorld, ",")
			if len(sliceOnlyWorld) < 1 {
				html, err := replaceParams(params, templateInfo.Content, templateInfo.FrameId)
				if err != nil {
					return nil, err
				}
				templates = append(templates, map[string]string{"title": templateInfo.Title, "html": html})
			}
		}
	}
	return templates, nil
}

func replaceParams(params map[string]string, content string, frameId int64) (html string, err error) {
	for k, v := range params {
		html = strings.Replace(content, "%"+k+"%", v, -1)
	}

	emailFrame, err := repo.EmailFrameRepo.GetInfoById(frameId)
	if err != nil {
		return "", err
	}
	if emailFrame == nil {
		return "", errors.New("not found id is" + strconv.Itoa(int(frameId)) + "email frame")
	}
	html = strings.Replace(emailFrame.Content, "%content%", html, -1)
	return html, nil
}

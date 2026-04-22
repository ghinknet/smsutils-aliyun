package aliyun

import (
	"strings"

	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	dysmsapi20180501 "github.com/alibabacloud-go/dysmsapi-20180501/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"go.gh.ink/smsutils/v3/errors"
	"go.gh.ink/smsutils/v3/model"
	"go.gh.ink/smsutils/v3/utils"
	"go.gh.ink/toolbox/pointer"
)

// sendMessageToGlobeRaw is the raw method to send an SMS to Globe
func sendMessageToGlobeRaw(c Client, dest string, sender string, message string) error {
	resp, err := c.Globe.SendMessageToGlobe(&dysmsapi20180501.SendMessageToGlobeRequest{
		To:      &dest,
		Message: &message,
		From:    &sender,
	})
	if err != nil {
		return err
	}
	if pointer.SafeDeref(resp.Body.ResponseCode) != "OK" {
		return errors.ErrDriverSendFailed.
			WithDriverName(Name).
			WithDriverCode(pointer.SafeDeref(resp.GetBody().GetResponseCode())).
			WithDriverMessage(pointer.SafeDeref(resp.GetBody().GetResponseDescription())).
			WithDriverRequestID(pointer.SafeDeref(resp.GetBody().GetRequestId())).
			WithDriverResponse(resp.GetBody())
	}

	return nil
}

// sendMessageToChineseMainlandRaw is the raw method to send an SMS to Chinese Mainland
func sendMessageToChineseMainlandRaw(c Client, dest string, sender string, template string, params string) error {
	resp, err := c.CN.SendSmsWithOptions(&dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  &dest,
		SignName:      &sender,
		TemplateCode:  &template,
		TemplateParam: &params,
	}, new(util.RuntimeOptions))
	if err != nil {
		return err
	}
	if pointer.SafeDeref(resp.Body.Code) != "OK" {
		return errors.ErrDriverSendFailed.
			WithDriverName(Name).
			WithDriverCode(pointer.SafeDeref(resp.GetBody().GetCode())).
			WithDriverMessage(pointer.SafeDeref(resp.GetBody().GetMessage())).
			WithDriverRequestID(pointer.SafeDeref(resp.GetBody().GetRequestId())).
			WithDriverResponse(resp.GetBody())
	}

	return nil
}

func (c Client) SendMessage(dest string, sender string, template string, vars model.Vars) error {
	// Try to parse number
	dest, _, _, regionCode, err := utils.ProcessNumberForChinese(dest)
	if err != nil {
		return err
	}

	// Chinese mainland
	if regionCode == "CN" {
		// Preprocess vars
		params := make(map[string]string)
		for _, v := range vars {
			params[v.Key] = v.Value
		}

		// Marshal params
		marshalled, err := c.Marshal(params)
		if err != nil {
			return err
		}

		// Send message
		return sendMessageToChineseMainlandRaw(c, dest, sender, template, string(marshalled))
	}

	// Globe

	// Render template
	templateContent, ok := c.GlobeTemplate[template]
	if !ok {
		return errors.ErrDriverSendFailed.
			WithDriverName(Name).
			WithDriverMessage("template not found")
	}
	for _, v := range vars {
		templateContent = strings.Replace(
			templateContent, strings.Join([]string{"${", v.Key, "}"}, ""), v.Value, -1,
		)
	}

	// Send message
	return sendMessageToGlobeRaw(c, strings.TrimPrefix(dest, "+"), sender, templateContent)
}

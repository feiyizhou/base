package clients

import (
	"errors"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	txErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

type TencentSMSConf struct {
	SecretID      string        `json:"secretID" mapstructure:"secretID"`
	SecretKey     string        `json:"secretKey" mapstructure:"secretKey"`
	SdkAppID      string        `json:"sdkAppID" mapstructure:"sdkAppID"`
	Host          string        `json:"host" mapstructure:"host"`
	Region        string        `json:"region" mapstructure:"region"`
	SmsVerifyCode smsVerifyCode `json:"smsVerifyCode" mapstructure:"smsVerifyCode"`
}

type smsVerifyCode struct {
	SignName   string `json:"signName" mapstructure:"signName"`
	TemplateID string `json:"templateID" mapstructure:"templateID"`
}

type TXSMSClient struct {
	Client     *sms.Client `json:"clients"`
	SignName   string      `json:"signName"`
	TemplateID string      `json:"templateID"`
	SdkAppID   string      `json:"sdkAppID"`
}

func NewTXSMSClient(conf TencentSMSConf) (*TXSMSClient, error) {
	// 实例化一个认证对象，入参需要传入腾讯云账户密钥对secretId,secretKey.
	credential := common.NewCredential(
		conf.SecretID,
		conf.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.Endpoint = conf.Host
	cpf.SignMethod = "HmacSHA1"
	client, err := sms.NewClient(credential, conf.Region, cpf)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Init Tencent sms clients failed, err: %v", err))
	}
	return &TXSMSClient{
		Client:     client,
		SignName:   conf.SmsVerifyCode.SignName,
		TemplateID: conf.SmsVerifyCode.TemplateID,
		SdkAppID:   conf.SdkAppID,
	}, err
}

func (sc *TXSMSClient) SendShortMessage(contentArr, phoneNums []string) (bool, error) {
	/* 实例化一个请求对象，根据调用的接口和实际情况*/
	request := sms.NewSendSmsRequest()
	// 应用 ID 可前往 [短信控制台](https://console.cloud.tencent.com/smsv2/app-manage) 查看
	request.SmsSdkAppId = common.StringPtr(sc.SdkAppID)
	// 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名
	request.SignName = common.StringPtr(sc.SignName)
	/* 模板 ID: 必须填写已审核通过的模板 ID */
	request.TemplateId = common.StringPtr(sc.TemplateID)
	/* 模板参数: 模板参数的个数需要与 TemplateId 对应模板的变量个数保持一致，若无模板参数，则设置为空*/
	request.TemplateParamSet = common.StringPtrs(contentArr)
	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs(phoneNums)
	// 通过client对象调用想要访问的接口，需要传入请求对象
	_, err := sc.Client.SendSms(request)
	// 处理异常
	if _, ok := err.(*txErrors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return false, err
	}
	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	if err != nil {
		return false, err
	}
	return true, nil
}

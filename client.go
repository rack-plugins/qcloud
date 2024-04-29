package qcloud

import (
	"github.com/fimreal/goutils/ezap"
	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func NewClient(endpoint string) *vpc.Client {
	sid := viper.GetString("qcloud.sid")
	skey := viper.GetString("qcloud.skey")
	region := viper.GetString("qcloud.region")
	ezap.Debugf("获取腾讯云连接配置，SECRET_ID: %s，SECRET_KEY: ***，REGION_ID: %s", sid, region)
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential := common.NewCredential(
		sid,
		skey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以默认
	cpf := profile.NewClientProfile()
	if viper.GetBool("qcloud.insecureskipverify") {
		cpf.UnsafeRetryOnConnectionFailure = true // 跳过证书验证
	}
	// cpf.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"
	cpf.HttpProfile.Endpoint = endpoint

	client, err := vpc.NewClient(credential, region, cpf)
	if err != nil {
		ezap.Error(err.Error())
	}
	return client
}

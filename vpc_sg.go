package qcloud

import (
	"net/http"

	"github.com/fimreal/goutils/ezap"
	"github.com/gin-gonic/gin"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func Allow(c *gin.Context) {
	// 检查传入 json 参数是否符合
	var sgrule SGRule
	// 新版本 gin c.ShouldBind 需要额外传入 json type 才能支持，这里用 ShouldBindJSON
	if err := c.ShouldBind(&sgrule); err != nil {
		ezap.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ezap.Debugf("请求添加安全组规则, 传入 IP 地址: %s, 安全组 ID: %s, 备注: %s, 策略: %s", sgrule.IP, sgrule.SGID, sgrule.Remark, sgrule.Policy)
	if !sgrule.verify() {
		ezap.Errorw("qcloud securitygroup cap", "desc", "传入参数不符合要求")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "传入参数不符合要求"})
		return
	}

	err := sgrule.authorize()
	if err != nil {
		ezap.Error("添加安全组规则出错: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ezap.Debug("成功添加安全组规则")
	c.JSON(http.StatusOK, gin.H{"result": "成功添加安全组规则"})
}

// authorize() 添加安全组规则
func (s *SGRule) authorize() error {
	c := NewClient("vpc.tencentcloudapi.com")

	request := vpc.NewCreateSecurityGroupPoliciesRequest()
	request.SecurityGroupId = &s.SGID
	request.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{
		Ingress: []*vpc.SecurityGroupPolicy{
			{
				// 导入位置，默认插入第一个位置
				PolicyIndex:       common.Int64Ptr(0),
				Action:            common.StringPtr(s.Policy),
				CidrBlock:         common.StringPtr(s.IP),
				Protocol:          common.StringPtr("ALL"),
				Port:              common.StringPtr("ALL"),
				PolicyDescription: common.StringPtr(s.Remark),
			},
		},
	}

	response, err := c.CreateSecurityGroupPolicies(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		ezap.Printf("An API error has returned: %s", err)
	}
	if err != nil {
		return err
	}
	// 输出json格式的字符串回包
	ezap.Debug("成功添加安全组规则: %s", response.ToJsonString())
	return nil
}

// verify() 检查传入参数是否正确
// 检查 IP 是否为内网 IP，如果是则返回错误
// 否则通过
func (s *SGRule) verify() bool {
	if !IsIPv4(s.IP) {
		ezap.Errorf("传入 ip[%s] 有误，仅支持 ipv4", s.IP)
		return false
	}
	if IsLanIPv4(s.IP) {
		ezap.Errorf("传入 ip[%s] 为内网 ip", s.IP)
		return false
	}
	return true
}

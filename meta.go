package qcloud

import (
	"github.com/fimreal/rack/module"
	"github.com/spf13/cobra"
)

const (
	ID            = "qcloud"
	Comment       = "qcloud api"
	RoutePrefix   = "/" + ID
	DefaultEnable = false
)

var Module = module.Module{
	ID:      ID,
	Comment: Comment,
	// gin route
	RouteFunc:   AddRoute,
	RoutePrefix: RoutePrefix,
	// cobra flag
	FlagFunc: ServeFlag,
}

func ServeFlag(serveCmd *cobra.Command) {
	serveCmd.Flags().Bool(ID, DefaultEnable, Comment)

	serveCmd.Flags().String("qcloud.sid", "", "SecretId")
	serveCmd.Flags().String("qcloud.skey", "", "SecretKey")
	serveCmd.Flags().String("qcloud.region", "", "region, eg. ap-beijing")
	serveCmd.Flags().Bool("qcloud.insecureskipverify", false, "是否跳过证书验证(小容器没有证书会遇到 https 连接证书验证失败)")
}

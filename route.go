package qcloud

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func AddRoute(g *gin.Engine) {
	if !viper.GetBool(ID) && !viper.GetBool("allservices") {
		return
	}
	g.GET("/help/qcloud", help)

	r := g.Group(RoutePrefix)
	r.POST("/addsgrule", Allow)
}

package router

import (
	"github.com/gin-gonic/gin"

	first "github.com/michaelchandrag/gobiz-spy/module/controller/first"
	second "github.com/michaelchandrag/gobiz-spy/module/controller/second"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("public/views/*")
	r.Static("/css", "public/assets/css")
	r.Static("/fonts", "public/assets/fonts")
	r.Static("/img", "public/assets/img")
	r.Static("/js", "public/assets/js")
	r.Static("/node_modules", "public/assets/node_modules")

	r.GET("/", first.RenderFirstPage)
	r.GET("/page2", second.RenderSecondPage)
	r.POST("/request_otp_gobiz", first.RequestOtpGobiz)
	r.POST("/request_token_gobiz", first.RequestTokenGobiz)
	r.POST("/request_profile_gobiz", first.RequestProfileGobiz)
	r.POST("/request_transactions_gobiz", second.RequestTransactionsGobiz)

	return r
}
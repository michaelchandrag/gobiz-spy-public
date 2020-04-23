package first

import (
	"os"
	"net/http"
	"fmt"

	gobiz "github.com/michaelchandrag/gobiz-spy/module/util/gobiz"
	"github.com/gin-gonic/gin"
)

func RenderFirstPage (c *gin.Context) {
	baseUrl := os.Getenv("BASE_URL")
	c.HTML(http.StatusOK, "page1.html", gin.H{
		"title": "Main website",
		"baseUrl": baseUrl,
	})
}

func RequestOtpGobiz (c *gin.Context) {
	type Body struct {
		PhoneNumber 	string 		`json:"phone_number"`
	}

	var payload Body
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H {
			"success": false,
			"message": "Body request is not correct.",
		})
		return
	}

	res, err := gobiz.RequestOtp(payload.PhoneNumber)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to request OTP.",
			"error": err,
			"response_data": res,
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Success request OTP. A SMS has been sent to your phone number.",
		"response_data": res,
	})
	return
}

func RequestTokenGobiz (c *gin.Context) {
	type Body struct {
		OtpToken 	string 			`json:"otp_token"`
		Otp 		string 			`json:"otp"`
	}

	var payload Body
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H {
			"success": false,
			"message": "Body request is not correct.",
		})
		return
	}

	res, err := gobiz.RequestToken(payload.OtpToken, payload.Otp)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to request token.",
			"error": err,
			"response_data": res,
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Success request token. You are redirected to the main menu option.",
		"response_data": res,
	})
	return
}

func RequestProfileGobiz (c *gin.Context) {
	type Body struct {
		AccessToken 	string 		`json:"access_token"`
	}

	var payload Body
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Body request is not correct.",
		})
	}

	res, err := gobiz.RequestProfile(payload.AccessToken)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to request profile merchant.",
			"error": err,
			"response_data": res,
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Success request profile merchant.",
		"response_data": res,
	})
	return
}
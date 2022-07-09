package controller

import (
	"eden-admin/app/admin/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginController struct {
}

func (l *LoginController) Captcha(c *gin.Context) {
	id, b64s, err := middleware.DriverDigitFunc()
	if err != nil {
		fmt.Printf("DriverDigitFunc error, %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    b64s,
		"id":      id,
		"message": "success",
	})

}

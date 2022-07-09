package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type R map[string]interface{}

func New(code int, message string) R {
	r := make(R, 3)
	r["code"] = code
	r["message"] = message
	return r
}

func Ok() R {
	return New(0, "成功")
}

func Err(code int, message string) R {
	return New(code, message)
}

func Errf(code int, template string) R {
	return nil
}

func (r R) Of(key string, value interface{}) R {
	r[key] = value
	return r
}

func (r R) Response(g *gin.Context) {
	g.JSON(http.StatusOK, r)
}

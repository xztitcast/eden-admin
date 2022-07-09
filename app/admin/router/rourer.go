package router

import (
	"eden-admin/app/admin/controller"
	"eden-admin/app/admin/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
)

var (
	routerNoCheckRole = make([]func(group *gin.RouterGroup), 0)
	routerCheckRole   = make([]func(g *gin.RouterGroup, jwtMiddleware *jwt.GinJWTMiddleware), 0)
)

func InitRouter() *gin.Engine {
	var r = gin.Default()
	g := r.Group("/sys")
	g.Use(gin.Logger()).Use(gin.Recovery())
	//g.Use(Recover)
	j := initJwtRouter(g)
	needNoRouterCheckRole(g)
	needRouterCheckRole(g, j)
	serializeInject()
	return r
}

func needRouterCheckRole(g *gin.RouterGroup, j *jwt.GinJWTMiddleware) {
	for _, f := range routerCheckRole {
		f(g, j)
	}
}

func needNoRouterCheckRole(g *gin.RouterGroup) {
	for _, f := range routerNoCheckRole {
		f(g)
	}
}

func initJwtRouter(e *gin.RouterGroup) *jwt.GinJWTMiddleware {
	var loginController controller.LoginController
	var jwt middleware.Jwt
	var middleware = jwt.GinJWTMMiddlewareInit()
	e.POST("/login", middleware.LoginHandler)
	e.GET("/refresh_token", middleware.RefreshHandler)
	e.GET("/captcha", loginController.Captcha)
	injectCollector = append(injectCollector, &inject.Object{Value: &jwt})
	return middleware
}

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			var errMessage string
			switch r.(type) {
			case runtime.Error:
				errMessage = r.(error).Error()
				log.Printf("runtime: %v\n", r)
			default:
				errMessage = "系统异常,稍后重试"
				log.Printf("panic: %v\n", r)
				debug.PrintStack()
			}
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  errMessage,
			})
			c.Abort()
		}
	}()
	c.Next()
}

func serializeInject() {
	var injector inject.Graph
	if err := injector.Provide(injectCollector...); err != nil {
		log.Fatal("inject fatal: ", err)
	}
	if err := injector.Populate(); err != nil {
		log.Fatal("injector fatal: ", err)
	}
}

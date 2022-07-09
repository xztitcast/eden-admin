package router

import (
	"eden-admin/app/admin/common"
	"eden-admin/app/admin/controller"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
)

func RegisterSysUserRouterAndInject(e *gin.RouterGroup, middleware *jwt.GinJWTMiddleware) {
	//注册路由
	var sysUserController controller.SysUserController
	u := e.Group("/user").Use(middleware.MiddlewareFunc()).Use(common.Subject())
	{
		u.GET("/info", sysUserController.Info)
		u.GET("/select/:userId", sysUserController.Select)
		u.GET("/list", sysUserController.List)
		u.POST("/save", sysUserController.Save)
		u.POST("/update", sysUserController.Update)
		u.POST("/delete", sysUserController.Delete)
	}

	//收集依赖
	injectCollector = append(injectCollector, &inject.Object{Value: &sysUserController})
}

func init() {
	routerCheckRole = append(routerCheckRole, RegisterSysUserRouterAndInject)
}

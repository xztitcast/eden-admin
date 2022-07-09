package router

import (
	"eden-admin/app/admin/controller"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
)

func RegisterSysMenuRouterAndInject(e *gin.RouterGroup, middleware *jwt.GinJWTMiddleware) {
	var sysMenuController controller.SysMenuController
	m := e.Group("/menu").Use(middleware.MiddlewareFunc())
	{
		m.GET("/nav", sysMenuController.Nav)
		m.GET("/list", sysMenuController.List)
		m.GET("/select", sysMenuController.Select)
		m.GET("/info", sysMenuController.Info)
		m.POST("/save", sysMenuController.Save)
		m.POST("/update", sysMenuController.Update)
		m.POST("/delete", sysMenuController.Delete)
	}
	injectCollector = append(injectCollector, &inject.Object{Value: &sysMenuController})
}

func init() {
	routerCheckRole = append(routerCheckRole, RegisterSysMenuRouterAndInject)
}

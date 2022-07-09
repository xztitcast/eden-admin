package router

import (
	"eden-admin/app/admin/common"
	"eden-admin/app/admin/controller"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
)

func RegisterSysRoleRouterAndInject(e *gin.RouterGroup, middleware *jwt.GinJWTMiddleware) {
	var sysRoleController controller.SysRoleController
	r := e.Group("/role").Use(middleware.MiddlewareFunc()).Use(common.Subject())
	{
		r.GET("/list", sysRoleController.List)
		r.GET("/select", sysRoleController.Select)
		r.GET("/info", sysRoleController.Info)
		r.POST("/save", sysRoleController.Save)
		r.POST("/update", sysRoleController.Update)
		r.POST("/delete", sysRoleController.Delete)
	}

	injectCollector = append(injectCollector, &inject.Object{Value: &sysRoleController})
}

func init() {
	routerCheckRole = append(routerCheckRole, RegisterSysRoleRouterAndInject)
}

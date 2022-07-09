package router

import (
	"eden-admin/app/admin/service"
	"eden-admin/logger"
	"eden-admin/mysql"
	"github.com/facebookgo/inject"
)

var injectCollector = make([]*inject.Object, 0)

func init() {
	injectCollector = append(injectCollector,
		&inject.Object{Value: mysql.NewDatabase()},
		&inject.Object{Value: logger.NewLogger()},
		&inject.Object{Value: &service.SysUserServiceImpl{}},
		&inject.Object{Value: &service.SysRoleServiceImpl{}},
		&inject.Object{Value: &service.SysMenuServiceImpl{}},
		&inject.Object{Value: &service.SysUserRoleServiceImpl{}},
		&inject.Object{Value: &service.SysRoleMenuServiceImpl{}},
	)
}

package controller

import (
	"eden-admin/app/admin/common"
	"eden-admin/app/admin/entity"
	"eden-admin/app/admin/service"
	"eden-admin/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type SysMenuController struct {
	Logger *logger.Logger         `inject:""`
	Sms    service.SysMenuService `inject:""`
	Sus    service.SysUserService `inject:""`
}

func (s *SysMenuController) Nav(g *gin.Context) {
	menuList := s.Sms.GetUserMenuList(1)
	_, permissions := s.Sus.GetAllPerms(1)
	common.Ok().Of("menuList", menuList).Of("permissions", permissions).Response(g)
}

func (s *SysMenuController) List(g *gin.Context) {
	menuList := s.Sms.GetMenuList()
	for i := 0; i < len(menuList); i++ {
		menu := s.Sms.GetById(menuList[i].ParentId)
		if menuList[i].ParentId == 0 {
			menuList[i].ParentName = "一级菜单"
		} else {
			menuList[i].ParentName = menu.Name
		}
	}
	g.JSON(http.StatusOK, menuList)
}

func (s *SysMenuController) Select(g *gin.Context) {
	menuList := s.Sms.GetNotButtonList()
	var root entity.SysMenu
	root.MenuId = 0
	root.Name = "一级菜单"
	root.ParentId = -1
	root.Open = true
	menuList = append(menuList, root)
	common.Ok().Of("menuList", menuList)
}

func (s *SysMenuController) Info(g *gin.Context) {
	menuId, _ := strconv.ParseInt(g.Query("menuId"), 10, 64)
	menu := s.Sms.GetById(menuId)
	common.Ok().Of("menu", menu).Response(g)
}

func (s *SysMenuController) Save(g *gin.Context) {
	var sysMenu entity.SysMenu
	if err := g.ShouldBind(&sysMenu); err != nil {
		s.Logger.Log.Errorf("SysMenuController.save()参数绑定异常 %s", err)
		common.Err(-1, "参数异常").Response(g)
		return
	}
	if err := s.verifyForm(&sysMenu); err != nil {
		s.Logger.Log.Warnf("SysMenuContoller.verifyForm()验证不通过 %s", err)
		common.Err(-1, err.Error()).Response(g)
		return
	}
	s.Sms.Save(&sysMenu)
	common.Ok().Response(g)
}

func (s *SysMenuController) Update(g *gin.Context) {
	var sysMenu entity.SysMenu
	if err := g.ShouldBind(&sysMenu); err != nil {
		s.Logger.Log.Errorf("SysMenuController.update()参数绑定异常 %s", err)
		common.Err(-1, "参数异常").Response(g)
		return
	}
	if err := s.verifyForm(&sysMenu); err != nil {
		s.Logger.Log.Warnf("SysMenuController.verifyForm()验证不通过 %s", err)
		common.Err(-1, err.Error()).Response(g)
		return
	}
	s.Sms.Update(&sysMenu)
	common.Ok().Response(g)
}

func (s *SysMenuController) Delete(g *gin.Context) {
	menuId, _ := strconv.ParseInt(g.Query("menuId"), 10, 64)
	if menuId <= 31 {
		common.Err(-1, "系统框架菜单不能被删除").Response(g)
		return
	}
	menuList := s.Sms.GetListParentIdf(menuId)
	if len(menuList) > 0 {
		common.Err(-1, "请删除对应的菜单或者按钮再删除主菜单").Response(g)
		return
	}
	s.Sms.Delete(menuId)
	common.Ok().Response(g)
}

func (s *SysMenuController) verifyForm(menu *entity.SysMenu) error {
	if len(strings.Trim(menu.Name, " ")) == 0 {
		return errors.New("菜单名称不能为空")
	}
	if menu.ParentId < 0 {
		return errors.New("上级菜单异常")
	}
	if menu.Type == service.MENU {
		if len(strings.Trim(menu.Url, " ")) == 0 {
			return errors.New("菜单URL不能为空")
		}
	}
	var parentType = service.CATALOG
	if menu.ParentId != 0 {
		sysMenu := s.Sms.GetById(menu.MenuId)
		parentType = int(sysMenu.Type)
	}

	if menu.Type == service.CATALOG || menu.Type == service.MENU {
		if parentType != service.CATALOG {
			return errors.New("上级菜单只能为目录")
		}
	}

	if menu.Type == service.BUTTON {
		if parentType != service.MENU {
			return errors.New("上级菜单只能为菜单")
		}
	}
	return nil
}

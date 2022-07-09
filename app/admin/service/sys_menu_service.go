package service

import (
	"eden-admin/app/admin/entity"
	"eden-admin/mysql"
	"gorm.io/gorm"
)

const (
	CATALOG = 0 //目录
	MENU    = 1 //菜单
	BUTTON  = 2 //按钮
)

type SysMenuService interface {
	GetById(menuId int64) *entity.SysMenu
	//根据父菜单获取子菜单(重载)
	GetListParentIdf(parentId int64) []entity.SysMenu
	//根据父菜单获取子菜单
	GetListParentId(parentId int64, menuIdList []int64) []entity.SysMenu
	//获取不包含按钮的菜单列表
	GetNotButtonList() []entity.SysMenu
	//获取用户菜单列表
	GetUserMenuList(userId int64) []entity.SysMenu

	GetMenuList() []entity.SysMenu

	Save(sysMenu *entity.SysMenu)

	Update(sysMenu *entity.SysMenu)

	//删除
	Delete(menuId int64)
}

type SysMenuServiceImpl struct {
	Db  *mysql.Database `inject:""`
	Sus SysUserService  `inject:""`
}

func (s *SysMenuServiceImpl) GetById(menuId int64) *entity.SysMenu {
	var sysMenu entity.SysMenu
	s.Db.Orm.Where("menu_id = ?", menuId).Find(&sysMenu)
	return &sysMenu
}

func (s *SysMenuServiceImpl) GetListParentIdf(parenId int64) []entity.SysMenu {
	var menuList = make([]entity.SysMenu, 0)
	s.Db.Orm.Where("parent_id = ?", parenId).Order("order_num ASC").Find(&menuList)
	return menuList
}

func (s *SysMenuServiceImpl) GetListParentId(parentId int64, menuIdList []int64) []entity.SysMenu {
	menuList := s.GetListParentIdf(parentId)
	if menuIdList == nil || len(menuIdList) == 0 {
		return menuList
	}
	list := make([]entity.SysMenu, 0)
	for _, v := range menuList {
		for _, m := range menuIdList {
			if v.MenuId == m {
				list = append(list, v)
			}
		}
	}
	return list
}

func (s *SysMenuServiceImpl) GetNotButtonList() []entity.SysMenu {
	menuList := make([]entity.SysMenu, 0)
	s.Db.Orm.Where("type != ?", 2).Order("order_num ASC").Find(&menuList)
	return menuList
}

func (s *SysMenuServiceImpl) GetUserMenuList(userId int64) []entity.SysMenu {
	if userId == 1 {
		return s.getAllMenuList(nil)
	}
	menuIdList := s.Sus.GetAllMenuId(userId)
	return s.getAllMenuList(menuIdList)
}

func (s *SysMenuServiceImpl) GetMenuList() []entity.SysMenu {
	menuList := make([]entity.SysMenu, 0)
	s.Db.Orm.Find(&menuList)
	return menuList
}

func (s *SysMenuServiceImpl) Save(sysMenu *entity.SysMenu) {
	s.Db.Orm.Save(sysMenu)
}

func (s *SysMenuServiceImpl) Update(sysMenu *entity.SysMenu) {
	s.Db.Orm.Updates(sysMenu)
}

func (s *SysMenuServiceImpl) Delete(menuId int64) {
	var sysMenu entity.SysMenu
	var sysRoleMenu entity.SysRoleMenu
	s.Db.Orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_id", menuId).Delete(&sysMenu).Error; err != nil {
			return err
		}
		if err := tx.Where("menu_id", menuId).Delete(&sysRoleMenu).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *SysMenuServiceImpl) getAllMenuList(menuIdList []int64) []entity.SysMenu {
	menuList := s.GetListParentId(0, menuIdList)
	s.getMenuTreeList(menuList, menuIdList)
	return menuList
}

func (s *SysMenuServiceImpl) getMenuTreeList(menuList []entity.SysMenu, menuIdList []int64) []entity.SysMenu {
	result := make([]entity.SysMenu, 0)
	for i := 0; i < len(menuList); i++ {
		if menuList[i].Type == CATALOG {
			menuList[i].List = s.getMenuTreeList(s.GetListParentId(menuList[i].MenuId, menuIdList), menuIdList)
		}
		result = append(result, menuList[i])
	}
	return result
}

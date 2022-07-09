package service

import (
	"eden-admin/app/admin/entity"
	"eden-admin/mysql"
	"gorm.io/gorm"
)

type SysRoleMenuService interface {
	//保存或更新
	SaveOrUpdate(roleId int64, menuIdList []int64, tx *gorm.DB) (err error)
	//获取菜列表
	GetMenuIdList(roleId int64) []int64
	//批量删除角色ID
	DeleteBatch(tx *gorm.DB, roleIds ...int64) (err error)
}

type SysRoleMenuServiceImpl struct {
	Db *mysql.Database `inject:""`
}

func (s *SysRoleMenuServiceImpl) SaveOrUpdate(roleId int64, menuIdList []int64, tx *gorm.DB) (err error) {
	if tx == nil {
		tx = s.Db.Orm
	}
	if err = s.DeleteBatch(tx, roleId); err != nil {
		return
	}
	if menuIdList == nil || len(menuIdList) == 0 {
		return nil
	}
	for _, v := range menuIdList {
		var sysRoleMenu entity.SysRoleMenu
		sysRoleMenu.RoleId = roleId
		sysRoleMenu.MenuId = v
		err = tx.Save(&sysRoleMenu).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SysRoleMenuServiceImpl) GetMenuIdList(roleId int64) []int64 {
	menuIdList := make([]int64, 0)
	s.Db.Orm.Where("role_id = ?", roleId).Scan(&menuIdList)
	return menuIdList
}

func (s *SysRoleMenuServiceImpl) DeleteBatch(tx *gorm.DB, roleIds ...int64) (err error) {
	var sysRoleMenu entity.SysRoleMenu
	if tx == nil {
		tx = s.Db.Orm
	}
	err = tx.Where("role_id IN ?", roleIds).Delete(&sysRoleMenu).Error
	return
}

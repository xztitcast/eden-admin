package service

import (
	"eden-admin/app/admin/entity"
	"eden-admin/mysql"
	"gorm.io/gorm"
)

type SysUserRoleService interface {
	SaveOrUpdate(tx *gorm.DB, userId int64, roleIdList ...int64) (err error)

	GetRoleIdList(userId int64) []int64

	DeleteBatch(tx *gorm.DB, roleIds ...int64) (err error)
}

type SysUserRoleServiceImpl struct {
	Db *mysql.Database `inject:""`
}

func (s *SysUserRoleServiceImpl) SaveOrUpdate(tx *gorm.DB, userId int64, roleIdList ...int64) (err error) {
	if tx == nil {
		tx = s.Db.Orm
	}
	var sysUserRole entity.SysUserRole
	if err = tx.Where("user_id = ?", userId).Delete(&sysUserRole).Error; err != nil {
		return
	}
	if roleIdList == nil || len(roleIdList) == 0 {
		return nil
	}
	for _, v := range roleIdList {
		sur := new(entity.SysUserRole)
		sur.RoleId = v
		sur.UserId = userId
		if err = tx.Save(sur).Error; err != nil {
			return
		}
	}
	return nil
}

func (s *SysUserRoleServiceImpl) GetRoleIdList(userId int64) []int64 {
	roleIdList := make([]int64, 0)
	s.Db.Orm.Where("user_id = ?", userId).Find(&roleIdList)
	return roleIdList
}

func (s *SysUserRoleServiceImpl) DeleteBatch(tx *gorm.DB, roleIds ...int64) (err error) {
	if tx == nil {
		tx = s.Db.Orm
	}
	var sysUserRole entity.SysUserRole
	err = tx.Where("role_id IN ?", roleIds).Delete(&sysUserRole).Error
	return
}

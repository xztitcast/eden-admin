package service

import (
	"eden-admin/app/admin/common"
	"eden-admin/app/admin/entity"
	"eden-admin/mysql"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type SysRoleService interface {
	//新增
	Save(sysRole *entity.SysRole)
	//更新
	Update(sysRole *entity.SysRole)
	//批量删除所有角色
	DeleteBatch(roleIds ...int64)
	//获取用户创建的角色ID列表
	GetRoleIdList(adder int64) []int64
	//获取角色列表
	GetSysRoleList(pageNum, pageSize int, roleName string, adder int64) *common.P

	GetByAdder(userId int64) []entity.SysRole

	GetById(roleId int64) *entity.SysRole
}

type SysRoleServiceImpl struct {
	Db   *mysql.Database    `inject:""`
	Sus  SysUserService     `inject:""`
	Srs  SysRoleMenuService `inject:""`
	Surs SysUserRoleService `inject:""`
}

func (s *SysRoleServiceImpl) Save(sysRole *entity.SysRole) {
	s.Db.Orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(sysRole).Error; err != nil {
			return err
		}
		if err := s.checkPerms(sysRole); err != nil {
			return err
		}
		err := s.Srs.SaveOrUpdate(sysRole.RoleId, sysRole.MenuIdList, tx)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *SysRoleServiceImpl) Update(sysRole *entity.SysRole) {
	s.Db.Orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(sysRole).Error; err != nil {
			return err
		}
		if err := s.checkPerms(sysRole); err != nil {
			return err
		}
		err := s.Srs.SaveOrUpdate(sysRole.RoleId, sysRole.MenuIdList, tx)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *SysRoleServiceImpl) DeleteBatch(roleIds ...int64) {
	var sysRole entity.SysRole
	s.Db.Orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&sysRole, roleIds).Error; err != nil {
			return err
		}
		if err := s.Srs.DeleteBatch(tx, roleIds...); err != nil {
			return err
		}
		if err := s.Surs.DeleteBatch(tx, roleIds...); err != nil {
			return err
		}
		return nil
	})
}

func (s *SysRoleServiceImpl) GetRoleIdList(adder int64) []int64 {
	roleIdList := make([]int64, 0)
	s.Db.Orm.Raw("SELECT role_id FROM sys_role WHERE adder = ?", adder).Find(&roleIdList)
	return roleIdList
}

func (s *SysRoleServiceImpl) GetSysRoleList(pageNum, pageSize int, roleName string, adder int64) *common.P {
	var count int64
	roleList := make([]entity.SysRole, 0)
	db := s.Db.Orm.Limit(pageSize).Offset((pageNum - 1) * pageSize)
	if len(strings.Trim(roleName, " ")) > 0 {
		db.Where("role_name = ?", roleName)
	}
	if adder > 0 {
		db.Where("adder = ?", adder)
	}
	db.Scan(&roleList).Count(&count)
	return common.NewP(count, roleList)
}

func (s *SysRoleServiceImpl) GetByAdder(userId int64) []entity.SysRole {
	list := make([]entity.SysRole, 0)
	s.Db.Orm.Table("sys_role").Where("adder = ?", userId).Find(&list)
	return list
}

func (s *SysRoleServiceImpl) GetById(roleId int64) *entity.SysRole {
	var sysRole entity.SysRole
	s.Db.Orm.Where("role_id = ?", roleId).Scan(&sysRole)
	return &sysRole
}

func (s *SysRoleServiceImpl) checkPerms(role *entity.SysRole) error {
	if role.Adder == 1 {
		return nil
	}
	menuIdList := s.Sus.GetAllMenuId(role.RoleId)
	flag := true
	for _, v := range role.MenuIdList {
		for _, m := range menuIdList {
			if v != m {
				flag = false
				goto Label
			}
		}
	}
Label:
	if !flag {
		return errors.New("新增角色的权限，已超出你的权限范围")
	}
	return nil
}

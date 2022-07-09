package service

import (
	"eden-admin/app/admin/common"
	"eden-admin/app/admin/entity"
	"eden-admin/mysql"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type SysUserService interface {
	Info(sysUser *entity.SysUser)
	//查询用户列表
	GetSysUserList(pageNum, pageSize int, username string, userId int64) *common.P
	//查询所有用户权限
	GetAllPerms(userId int64) (map[string]bool, []string)
	//查询用户的所有菜单ID
	GetAllMenuId(userId int64) []int64
	//根据用户名,查询系统用户
	GetByUsername(username string) *entity.SysUser
	//删除用户
	DeleteBatch(userIds ...int64)
	//保存
	Save(sysUser *entity.SysUser) error
	//更新
	Update(sysUser *entity.SysUser) error
}

type SysUserServiceImpl struct {
	Db   *mysql.Database    `inject:""`
	Surs SysUserRoleService `inject:""`
}

func (s *SysUserServiceImpl) Info(sysUser *entity.SysUser) {
	s.Db.Orm.Table("sys_user").Select([]string{"user_id", "username", "status", "avatar", "adder", "created"}).Find(sysUser)
}

func (s *SysUserServiceImpl) GetSysUserList(pageNum, pageSize int, username string, userId int64) *common.P {
	var count int64
	sysUsers := make([]entity.SysUser, 0)
	tx := s.Db.Orm.Limit(pageSize).Offset((pageNum - 1) * pageSize)
	if len(strings.Trim(username, " ")) > 0 {
		tx.Where("username = ?", username)
	}
	if userId != 1 {
		tx.Where("adder = ?", userId)
	}
	tx.Select([]string{"user_id", "username", "status", "avatar", "adder", "created"}).Find(&sysUsers).Count(&count)
	return common.NewP(count, sysUsers)
}

func (s *SysUserServiceImpl) GetAllPerms(userId int64) (map[string]bool, []string) {
	perms := make([]string, 0)
	if userId == 1 {
		s.Db.Orm.Raw("SELECT perms FROM sys_menu").Scan(&perms)
	} else {
		s.Db.Orm.Raw("SELECT m.perms FROM sys_user_role ur LEFT JOIN sys_role_menu rm ON ur.role_id = rm.role_id LEFT JOIN sys_menu m ON rm.menu_id = m.menu_id WHERE ur.user_id = ?", userId).Scan(&perms)
	}
	dataList := make([]string, 0)
	dataMap := make(map[string]bool, 0)
	for _, v := range perms {
		if len(strings.Trim(v, " ")) == 0 {
			continue
		}
		split := strings.Split(v, ",") //截取
		for i := 0; i < len(split); i++ {
			if _, ok := dataMap[split[i]]; !ok { //利用map去重
				dataMap[split[i]] = true
				dataList = append(dataList, split[i])
			}
		}
	}
	return dataMap, dataList
}

func (s *SysUserServiceImpl) GetAllMenuId(userId int64) []int64 {
	menuIdList := make([]int64, 0)
	s.Db.Orm.Raw("SELECT distinct rm.menu_id FROM sys_user_role ur LEFT JOIN sys_role_menu rm ON ur.role_id = rm.role_id WHERE ur.user_id = ?", userId).Scan(&menuIdList)
	return menuIdList
}

func (s *SysUserServiceImpl) GetByUsername(username string) *entity.SysUser {
	var sysUser entity.SysUser
	s.Db.Orm.Where("username = ?", username).Find(&sysUser)
	return &sysUser
}

func (s *SysUserServiceImpl) DeleteBatch(userIds ...int64) {
	var sysUser entity.SysUser
	s.Db.Orm.Delete(&sysUser, userIds)
}

func (s *SysUserServiceImpl) Save(sysUser *entity.SysUser) error {
	return s.Db.Orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(sysUser).Error; err != nil {
			return err
		}
		if err := s.checkRole(sysUser); err != nil {
			return err
		}
		if err := s.Surs.SaveOrUpdate(tx, sysUser.UserId, sysUser.RoleIdList...); err != nil {
			return err
		}
		return nil
	})
}

func (s *SysUserServiceImpl) Update(sysUser *entity.SysUser) error {
	return s.Db.Orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(sysUser).Updates(sysUser).Error; err != nil {
			return err
		}
		if err := s.checkRole(sysUser); err != nil {
			return err
		}
		if err := s.Surs.SaveOrUpdate(tx, sysUser.UserId, sysUser.RoleIdList...); err != nil {
			return err
		}
		return nil
	})
}

func (s *SysUserServiceImpl) checkRole(sysUser *entity.SysUser) error {
	roleIdList := sysUser.RoleIdList
	if roleIdList == nil || len(roleIdList) == 0 {
		return nil
	}
	if sysUser.Adder == 1 {
		return nil
	}
	idList := s.Surs.GetRoleIdList(sysUser.UserId)
	flag := true
	for _, v := range roleIdList {
		for _, m := range idList {
			if v != m {
				flag = false
				goto Label
			}
		}
	}
Label:
	if !flag {
		return errors.New("新增用户所选角色,不是本人创建")
	}
	return nil
}

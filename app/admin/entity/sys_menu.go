package entity

import "reflect"

type SysMenu struct {
	MenuId     int64     `gorm:"primaryKey;autoIncrement;comment:主键菜单ID" json:"menuId" form:"menuId"`
	ParentId   int64     `gorm:"comment:父菜单ID,一级菜单为0" json:"parentId" form:"parentId"`
	ParentName string    `gorm:"-" json:"parentName" form:"parentName"`
	Name       string    `gorm:"size:50;comment:菜单名称" json:"name" form:"name"`
	Url        string    `gorm:"size:200;comment:菜单url" json:"url" form:"url"`
	Perms      string    `gorm:"size:500;comment:权限" json:"perms" form:"perms"`
	Type       int32     `gorm:"comment:类型 0:目录 1:菜单 2:按钮" json:"type" form:"type"`
	Icon       string    `gorm:"size:50;comment:图标" json:"icon" form:"icon"`
	OrderNum   int32     `gorm:"comment:排序" json:"orderNum" form:"orderNum"`
	Open       bool      `gorm:"-" json:"open" form:"open"`
	List       []SysMenu `gorm:"-" json:"list" form:"list"`
}

func (s SysMenu) IsEmpty() bool {
	return reflect.DeepEqual(s, SysMenu{})
}

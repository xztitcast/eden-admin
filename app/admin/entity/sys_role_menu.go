package entity

type SysRoleMenu struct {
	Id     int64 `gorm:"not null;primaryKey;autoIncrement;comment:主键ID" json:"id"`
	RoleId int64 `gorm:"comment:角色ID" json:"roleId"`
	MenuId int64 `gorm:"comment:菜单ID" json:"menuId"`
}

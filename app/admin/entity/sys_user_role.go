package entity

type SysUserRole struct {
	Id     int64 `gorm:"not null;primaryKey;autoIncrement;comment:主键ID" json:"id"`
	UserId int64 `gorm:"comment:用户ID" json:"userId"`
	RoleId int64 `gorm:"comment:角色ID" json:"roleId"`
}

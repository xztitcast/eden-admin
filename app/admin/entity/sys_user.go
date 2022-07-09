package entity

type SysUser struct {
	UserId     int64   `gorm:"primaryKey;autoIncrement;comment:编码" json:"userId" form:"userId"`
	Username   string  `gorm:"size:64;comment:用户名" json:"username" form:"username" binding:"required"`
	Password   string  `gorm:"size:128;comment:密码" json:"password" form:"password"`
	Salt       string  `gorm:"size:255;comment:盐" json:"salt" form:"salt"`
	Status     bool    `gorm:"size:1;comment:状态" json:"status" form:"status" binding:"required"`
	Adder      int64   `gorm:"comment:创建者ID" json:"adder" form:"adder"`
	Avatar     string  `gorm:"size:255;comment:头像" json:"avatar" form:"avatar"`
	RoleIdList []int64 `gorm:"-" json:"roleIdList" from:"roleIdList"`
	BaseEntity
}

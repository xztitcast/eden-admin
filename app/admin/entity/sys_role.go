package entity

type SysRole struct {
	RoleId     int64   `gorm:"not null;primaryKey;autoIncrement;comment:角色ID" json:"roleId" form:"roleId"`
	RoleName   string  `gorm:"size:100;comment:角色名称" json:"roleName" form:"roleName"`
	Remark     string  `gorm:"size:100;comment:备注" json:"remark" form:"remark"`
	Adder      int64   `gorm:"comment:创建者ID" json:"adder" form:"adder"`
	MenuIdList []int64 `gorm:"-" json:"menuIdList" form:"menuIdList"`
	BaseEntity
}

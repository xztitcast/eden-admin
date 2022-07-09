package vo

type UserPermVo struct {
	Username string   `json:"username"`
	Perms    []string `json:"perms"`
}

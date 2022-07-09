package middleware

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Captcha  string `form:"captcha" json:"captcha" binding:"required"`
	UUID     string `form:"uuid" json:"uuid" binding:"required"`
}

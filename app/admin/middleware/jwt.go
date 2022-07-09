package middleware

import (
	"eden-admin/app/admin/common/utils"
	"eden-admin/app/admin/entity"
	"eden-admin/app/admin/service"
	"errors"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type Jwt struct {
	Log log.Logger
	Sus service.SysUserService `inject:""`
}

func (j *Jwt) GinJWTMMiddlewareInit() (authMiddleware *jwt.GinJWTMiddleware) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Minute * 60,
		MaxRefresh:  time.Hour,
		IdentityKey: "idname",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(map[string]interface{}); ok {
				user, _ := v["user"]
				perms, _ := v["perms"]
				return jwt.MapClaims{
					"user":  user,
					"perms": perms,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return map[string]interface{}{
				"user":  claims["user"],
				"perms": claims["perms"],
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login Login
			if err := c.ShouldBind(&login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			if !Verify(login.UUID, login.Captcha, true) {
				return nil, errors.New("验证码错误")
			}
			user := j.Sus.GetByUsername(login.Username)
			if user.Password == utils.Sha256(login.Password, user.Salt) {
				permMap, _ := j.Sus.GetAllPerms(user.UserId)
				return map[string]interface{}{
					"user": &entity.SysUser{
						UserId:     user.UserId,
						Username:   user.Username,
						Status:     user.Status,
						Adder:      user.Adder,
						Avatar:     user.Avatar,
						RoleIdList: user.RoleIdList,
					},
					"perms": permMap,
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(map[string]interface{}); ok {
				return HasPerms(v["perms"].(map[string]interface{}), c)
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		j.Log.Fatal("JWT Error:" + err.Error())
	}
	return
}

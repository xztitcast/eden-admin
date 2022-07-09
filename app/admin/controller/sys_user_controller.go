package controller

import (
	"eden-admin/app/admin/common"
	"eden-admin/app/admin/common/utils"
	"eden-admin/app/admin/entity"
	"eden-admin/app/admin/service"
	"eden-admin/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

type SysUserController struct {
	Logger *logger.Logger         `inject:""`
	Sus    service.SysUserService `inject:""`
}

func (s *SysUserController) Info(g *gin.Context) {
	if subject, ok := g.Get("subject"); ok {
		common.Ok().Of("user", subject).Response(g)
	} else {
		common.Err(-1, "获取用户信息失败")
	}
}

func (s *SysUserController) Select(g *gin.Context) {
	var sysUser entity.SysUser
	s.Sus.Info(&sysUser)
	common.Ok().Of("result", sysUser).Response(g)
}

func (s *SysUserController) List(g *gin.Context) {
	var idx common.Index
	subject, ok := g.Get("subject")
	if !ok {
		common.Err(-1, "您未登录请重新登陆").Response(g)
		return
	}
	if err := g.ShouldBind(&idx); err != nil {
		common.Err(-1, "分页参数有误,请重新填写").Response(g)
		return
	}
	user := subject.(map[string]interface{})
	userId, err := strconv.ParseInt(fmt.Sprintf("%1.0f", user["userId"]), 10, 64)
	if err != nil {
		common.Err(-1, err.Error()).Response(g)
	}
	p := s.Sus.GetSysUserList(idx.PageNum, idx.PageSize, g.Query("username"), userId)
	common.Ok().Of("result", p).Response(g)

}

func (s *SysUserController) Save(g *gin.Context) {
	var sysUser entity.SysUser
	if err := g.ShouldBind(&sysUser); err != nil {
		common.Err(-1, "参数有误!").Response(g)
		s.Logger.Log.Errorf("SysUserController.save()绑定参数异常%s", err.Error())
		return
	}
	sysUser.Created = entity.JSONTime(time.Now())
	subject, _ := g.Get("subject")
	user := subject.(map[string]interface{})
	userId, _ := strconv.ParseInt(fmt.Sprintf("%1.0f", user["userId"]), 10, 64)
	sysUser.Adder = userId
	sysUser.Salt = uuid.New().String()
	sysUser.Password = utils.Sha256(user["password"].(string), sysUser.Salt)
	s.Sus.Save(&sysUser)
}

func (s *SysUserController) Update(g *gin.Context) {
	var sysUser entity.SysUser
	if err := g.ShouldBind(&sysUser); err != nil {
		common.Err(-1, "参数有误").Response(g)
		return
	}
	password := sysUser.Password
	if len(strings.Trim(password, " ")) != 0 {
		sysUser.Salt = uuid.New().String()
		sysUser.Password = utils.Sha256(password, sysUser.Salt)
	}
	subject, _ := g.Get("subject")
	user := subject.(map[string]interface{})
	userId, _ := strconv.ParseInt(fmt.Sprintf("%1.0f", user["userId"]), 10, 64)
	sysUser.Adder = userId
	s.Sus.Update(&sysUser)
}

func (s *SysUserController) Delete(g *gin.Context) {
	ids := make([]int64, 0)
	userIds := g.QueryArray("userIds")
	for _, v := range userIds {
		userId, _ := strconv.ParseInt(v, 10, 64)
		ids = append(ids, userId)
	}
	s.Sus.DeleteBatch(ids...)
}

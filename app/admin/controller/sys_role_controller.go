package controller

import (
	"eden-admin/app/admin/common"
	"eden-admin/app/admin/entity"
	"eden-admin/app/admin/service"
	"eden-admin/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type SysRoleController struct {
	Logger *logger.Logger             `inject:""`
	Srs    service.SysRoleService     `inject:""`
	Srms   service.SysRoleMenuService `inject:""`
}

func (s *SysRoleController) List(g *gin.Context) {
	var idx common.Index
	if err := g.ShouldBind(&idx); err != nil {
		s.Logger.Log.Errorf("SysRoleController.List() 分页参数异常 %s", err)
		common.Err(-1, "分页参数异常").Response(g)
		return
	}
	subject, _ := g.Get("subject")
	user := subject.(map[string]interface{})
	userId, _ := strconv.ParseInt(fmt.Sprintf("%1.0f", user["userId"]), 10, 64)
	if userId == 1 {
		userId = 0
	}
	p := s.Srs.GetSysRoleList(idx.PageNum, idx.PageSize, g.Query("roleName"), userId)
	common.Ok().Of("result", p).Response(g)
}

func (s *SysRoleController) Select(g *gin.Context) {
	subject, _ := g.Get("subject")
	user := subject.(map[string]interface{})
	userId, _ := strconv.ParseInt(fmt.Sprintf("%1.0f", user["userId"]), 10, 64)
	list := s.Srs.GetByAdder(userId)
	common.Ok().Of("result", list).Response(g)
}

func (s *SysRoleController) Info(g *gin.Context) {
	roleId, _ := strconv.ParseInt(g.Query("roleId"), 10, 64)
	sysRole := s.Srs.GetById(roleId)
	common.Ok().Of("result", sysRole)
}

func (s *SysRoleController) Save(g *gin.Context) {
	subject, _ := g.Get("subject")
	user := subject.(map[string]interface{})
	userId, _ := strconv.ParseInt(fmt.Sprintf("%1.0f", user["userId"]), 10, 64)
	var sysRole entity.SysRole
	if err := g.ShouldBind(&sysRole); err != nil {
		s.Logger.Log.Errorf("SysRoleController.save()参数绑定错误: %s", err)
		common.Err(-1, "参数错误!").Response(g)
		return
	}
	sysRole.Adder = userId
	sysRole.Created = entity.JSONTime(time.Now())
	s.Srs.Save(&sysRole)
	common.Ok().Response(g)
}

func (s *SysRoleController) Update(g *gin.Context) {
	subject, _ := g.Get("subject")
	user := subject.(map[string]interface{})
	userId, _ := strconv.ParseInt(fmt.Sprintf("%1.0f", user["userId"]), 10, 64)
	var sysRole entity.SysRole
	if err := g.ShouldBind(&sysRole); err != nil {
		s.Logger.Log.Errorf("SysRoleController.update()参数绑定错误: %s", err)
		common.Err(-1, "参数错误").Response(g)
		return
	}
	sysRole.Adder = userId
	s.Srs.Update(&sysRole)
	common.Ok().Response(g)
}

func (s *SysRoleController) Delete(g *gin.Context) {
	ids := make([]int64, 0)
	roleIds := g.QueryArray("roleIds")
	for _, v := range roleIds {
		id, _ := strconv.ParseInt(v, 10, 64)
		ids = append(ids, id)
	}
	s.Srs.DeleteBatch(ids...)
}

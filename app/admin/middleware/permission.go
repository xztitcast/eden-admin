package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	AND = "&"
	OR  = "|"
)

/**
请求权限配置
BASE 是一个map，存储的是请求路径与权限的映射关系，其key对应的请求路径，value对应的是请求权限
当权限中出现&符号表示同时满足多种权限，出现|符号表示满足其中一种即可
REST 是一个map，和BASE差不多只是存储的REST FUL风格的请求路径
*/
var (
	BASE = map[string]string{
		"/sys/user/info":   "sys:user:info",
		"/sys/user/list":   "sys:user:list",
		"/sys/user/save":   "sys:user:save",
		"/sys/user/update": "sys:user:update",
		"/sys/user/delete": "sys:user:delete",
		"/sys/role/list":   "sys:role:list",
		"/sys/role/select": "sys:role:select",
		"/sys/role/info":   "sys:role:info",
		"/sys/role/save":   "sys:role:save",
		"/sys/role/update": "sys:role:update",
		"/sys/role/delete": "sys:role:delete",
		"/sys/menu/list":   "sys:menu:list",
		"/sys/menu/select": "sys:menu:select",
		"/sys/menu/info":   "sys:menu:info",
		"/sys/menu/save":   "sys:menu:save",
		"/sys/menu/update": "sys:menu:update",
		"/sys/menu/delete": "sys:menu:delete",
	}
	REST = map[string]string{
		"/sys/user/select/:userId": "sys:user:select",
	}
)

func HasPerms(perms map[string]interface{}, g *gin.Context) bool {
	var perm string
	path := g.Request.URL.Path
	if len(g.Params) > 0 {
		for k, v := range REST {
			keys := strings.Split(k, ":")
			var builder strings.Builder
			str := keys[0]
			idx := strings.LastIndex(str, "/")
			builder.WriteString(str[:idx])
			data := keys[1:]
			for _, m := range data {
				builder.WriteString("/")
				builder.WriteString(g.Param(m))
			}
			if path == builder.String() {
				perm = v
				break
			}
		}
	} else {
		perm = BASE[path]
	}
	if len(perm) == 0 {
		return true
	}
	if strings.Contains(perm, AND) {
		temp := strings.Split(perm, AND)
		for _, v := range temp {
			if _, ok := perms[v]; !ok {
				return false
			}
		}
		return true
	}
	if strings.Contains(perm, OR) {
		temp := strings.Split(perm, OR)
		for _, v := range temp {
			if _, ok := perms[v]; ok {
				return true
			}
		}
		return false
	}
	_, ok := perms[perm]
	return ok
}

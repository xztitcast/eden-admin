package controller

import (
	jwt "github.com/appleboy/gin-jwt/v2"
)

type BaseController struct {
	Middleware *jwt.GinJWTMiddleware `inject:""`
}

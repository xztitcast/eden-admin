package common

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func Subject() gin.HandlerFunc {
	return func(g *gin.Context) {
		if data, ok := g.Get("JWT_PAYLOAD"); ok {
			result := data.(jwt.MapClaims)
			subject := result["user"].(map[string]interface{})
			g.Set("subject", subject)
		}

	}
}

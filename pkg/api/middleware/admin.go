package middlewares

import (
	"net/http"
	"strings"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AdminMiddleware struct{}

func NewAdminMiddleware() *AdminMiddleware {
	return &AdminMiddleware{}
}

func (mid *AdminMiddleware) AdminAuthMiddleware(c *gin.Context) {
	tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	role, err := utils.GetRoleFromToken(tokenString)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	if role == "admin" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

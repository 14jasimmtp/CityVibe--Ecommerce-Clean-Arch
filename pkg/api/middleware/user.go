package middlewares

import (
	"net/http"
	"strings"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserMiddleware struct{}

func NewUserMiddleware() *UserMiddleware {
	return &UserMiddleware{}
}

func (mid *UserMiddleware) UserAuthMiddleware(c *gin.Context) {
	tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	role, err := utils.GetRoleFromToken(tokenString)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	if role == "user" {
		c.Next()

	}
}

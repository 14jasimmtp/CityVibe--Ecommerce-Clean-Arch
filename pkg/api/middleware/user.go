package middlewares

import (
	"net/http"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserMiddleware struct{}

func NewUserMiddleware() *UserMiddleware{
	return &UserMiddleware{}
}

func (mid *UserMiddleware) UserAuthMiddleware(c *gin.Context) {
	Token, err := c.Cookie("Authorisation")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login to view page"})
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	role, err := utils.GetRoleFromToken(Token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	if role == "user" {
		c.Next()

	}
}

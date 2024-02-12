package middlewares

import (
	"net/http"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
	"github.com/gin-gonic/gin"
)


type AdminMiddleware struct{}

func NewAdminMiddleware()*AdminMiddleware{
	return &AdminMiddleware{}
}

func (mid *AdminMiddleware)AdminAuthMiddleware(c *gin.Context) {
	Token, err := c.Cookie("Authorisation")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	role, err := utils.GetRoleFromToken(Token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	if role == "admin" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

package middleware

import (
	"github.com/gin-gonic/gin"
)

type AuthNProvider interface {
	Authn(c *gin.Context)
}

type authProvider struct {
}

func NewAuthProvider() AuthNProvider {
	return &authProvider{}
}

func (ap *authProvider) Authn(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}
	c.Next()
}

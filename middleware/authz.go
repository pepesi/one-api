package middleware

import "github.com/gin-gonic/gin"

type AuthzProvider interface {
	Authz(c *gin.Context)
}

type authzProvider struct {
}

func NewAuthzProvider() AuthzProvider {
	return &authzProvider{}
}

func (ap *authzProvider) Authz(c *gin.Context) {
	c.Next()
}

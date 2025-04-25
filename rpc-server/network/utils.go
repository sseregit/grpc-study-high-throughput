package network

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (n *Network) verifyLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := getAuthToken(c)

		if t == "" {
			c.JSON(http.StatusUnauthorized, nil)
			c.Abort()
		} else {
			if _, err := n.gRPCClient.VerifyAuth(t); err != nil {
				c.JSON(http.StatusUnauthorized, err.Error())
				c.Abort()
			} else {
				c.Next()
			}
		}

	}
}

func getAuthToken(c *gin.Context) string {
	var token string

	authToken := c.Request.Header.Get("Authorization")
	authSliced := strings.Split(authToken, " ")
	if len(authSliced) > 1 {
		token = authSliced[1]
	}

	return token
}

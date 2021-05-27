package rest

import (
	"github.com/gin-gonic/gin"
)

const CredentialContextKey = "credentials"

const headerCredentialsKey = "Plooral-Credentials"

func ExtractCredentials() gin.HandlerFunc {
	return func(c *gin.Context) {
		credentials := c.Request.Header[headerCredentialsKey]
		if credentials != nil && len(credentials) == 1{
			c.Set(CredentialContextKey, credentials[0])
		}
	}
}
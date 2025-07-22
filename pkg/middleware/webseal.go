package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func Webseal(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. check that request comse from Webseal

		// 2. extract user id
		websealUserId := c.GetHeader("X-Webseal-User-ID")
		if websealUserId == "" {
			log.Error("Missing X-Webseal-User-ID header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			return
		}

		// 3. verify Webseal session
		//isValid, err := validateWebsealSession(c.Request) // todo: extract and call webseal auth endpoint
		//if err != nil || !isValid {
		//	log.Error("Invalid Webseal session", zap.Error(err))
		//	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		//	return
		//}

		// 4. query ldap for ougroups
		websealOuGroups := []string{"ou_group1", "ou_group2"} // todo: this should be fetched from LDAP
		//websealOuGroups, err := ldapClient.GetUserGroups(websealUserId)
		//if err != nil {
		//	log.Error("Failed to fetch LDAP groups", zap.Error(err))
		//	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "LDAP lookup failed"})
		//	return
		//}

		// inject user and ou groups into Gin context
		c.Set("user_id", websealUserId)
		c.Set("ou_groups", websealOuGroups)
		c.Next()
	}
}

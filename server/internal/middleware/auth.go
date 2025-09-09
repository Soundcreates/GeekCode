package middleware

import (
	"fmt"
	"net/http"
	"strings" 	
	"github.com/gin-gonic/gin"

	"geekCode/internal/auth"
)


func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) { //anonymous function
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("Auth Header: %s\n", authHeader) // Debugging line to print the auth header
		//checking if auth header is present
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error" : "missing auth header"})
			return
		}

		parts := strings.Split(authHeader, " ") 
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Printf("Invalid Auth Header Format: %s\n", authHeader) // Debugging line to print the invalid format
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error" : "invalid auth header"})
			return
		}
		fmt.Printf("Token: %s\n", parts[1]) // Debugging line to print the token part	
		userID, err := auth.ValidateToken(parts[1], secret)

		if err != nil {
			fmt.Printf("Token Validation Error: %v\n", err) // Debugging line to print the error
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error" : "invalid auth token"})
			return
		}


		//attaching user id to context (works like localstorage in js) , but its the gins context storage
		/* 
		Lives only for the current HTTP request being processed.

		It attaches data (like userId) to Gin’s context, so later middleware/handlers can retrieve it using c.Get("userId").

		Once the request finishes, the data is gone — it doesn’t persist across requests.
		*/

		c.Set("userId", userID)
		c.Next() 
	}
}

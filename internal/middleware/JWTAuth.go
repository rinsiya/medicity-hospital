package middleware

import (
    "net/http"
    "medicity/pkg/utils"

    "github.com/gin-gonic/gin"
)
func JWTAuth(role string) gin.HandlerFunc {

    return func(c *gin.Context) {

        token, err := c.Cookie("access_token")

        if err != nil || token == "" {
            c.Redirect(
                http.StatusSeeOther,
                "/"+role+"/login",
            )
            c.Abort()
            return
        }

        // Validate JWT first
        claims, err := utils.ValidateJWT(token)

        if err != nil {
            // Invalid/expired token
            c.SetCookie(
                "access_token",
                "",
                -1,
                "/",
                "",
                false,
                true,
            )

            c.Redirect(
                http.StatusSeeOther,
                "/"+role+"/login",
            )

            c.Abort()
            return
        }

        // JWT is valid, so now extract the values
        c.Set("userID", claims.UserID)
        c.Set("RoleID", claims.RoleID)
        c.Set("role", claims.Role)

        c.Next()
    }
}
func RequireRole(requiredRole string) gin.HandlerFunc {

    return func(c *gin.Context) {

        role, exists := c.Get("role")

        if !exists {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        if role.(string) != requiredRole {

            c.AbortWithStatus(http.StatusForbidden)
            return
        }

        c.Next()
    }
}
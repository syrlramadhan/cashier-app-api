package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/syrlramadhan/cashier-app/dto"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Authorization header is required",
			})
			return
		}

		// Check Bearer prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Invalid authorization header format. Use: Bearer <token>",
			})
			return
		}

		tokenString := tokenParts[1]
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "your-secret-key"
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Invalid token: " + err.Error(),
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract user info from claims with nil checks
			userIDClaim, exists := claims["user_id"]
			if !exists || userIDClaim == nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
					Success: false,
					Message: "Invalid token: user_id not found",
				})
				return
			}

			userIDFloat, ok := userIDClaim.(float64)
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
					Success: false,
					Message: "Invalid token: invalid user_id format",
				})
				return
			}

			emailClaim, exists := claims["email"]
			if !exists || emailClaim == nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
					Success: false,
					Message: "Invalid token: email not found",
				})
				return
			}

			email, ok := emailClaim.(string)
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
					Success: false,
					Message: "Invalid token: invalid email format",
				})
				return
			}

			roleClaim, exists := claims["role"]
			if !exists || roleClaim == nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
					Success: false,
					Message: "Invalid token: role not found",
				})
				return
			}

			role, ok := roleClaim.(string)
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
					Success: false,
					Message: "Invalid token: invalid role format",
				})
				return
			}

			userID := uint(userIDFloat)

			// Set user info in context
			ctx.Set("userID", userID)
			ctx.Set("email", email)
			ctx.Set("role", role)

			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Invalid token claims",
			})
			return
		}
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dto.APIResponse{
				Success: false,
				Message: "Access denied",
			})
			return
		}
		userRole := role.(string)
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, dto.APIResponse{
			Success: false,
			Message: "Insufficient permissions",
		})
	}
}

// AdminOnly middleware - only allows admin users
func AdminOnly() gin.HandlerFunc {
	return RoleMiddleware("admin")
}

// ManagerOrAdmin middleware - allows manager or admin users
func ManagerOrAdmin() gin.HandlerFunc {
	return RoleMiddleware("admin", "manager")
}

package jwt

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserRole string `json:"user_role"`
	jwt.RegisteredClaims
}

var AvailableUserRoles = []string{"employee", "moderator"}

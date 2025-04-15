package accesscontrol

import (
	"errors"
	"slices"

	"github.com/golang-jwt/jwt/v5"

	my_jwt "github.com/maximmihin/as25/internal/services/accesscontrol/pkg/jwt" // TODO rename alias
)

var ErrInvalidRole = errors.New("invalid roles") // TODO more info in message?

func NewDummyJWT(jwtSecret, role string) (string, error) {

	if !slices.Contains(my_jwt.AvailableUserRoles, role) {
		return "", ErrInvalidRole
	}

	claims := my_jwt.CustomClaims{
		UserRole: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}

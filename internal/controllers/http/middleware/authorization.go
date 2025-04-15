package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	accesscontrol "github.com/maximmihin/as25/internal/services/accesscontrol/pkg/jwt"
)

type AccessSet byte

func (as AccessSet) Has(set AccessSet) bool {
	return as&set != 0
}

const (
	AccessModerator AccessSet = 1 << iota
	AccessEmployee
)

var ErrAccessDenied = fiber.NewError(403)
var ErrUnexpectedUserRole = fiber.NewError(403, "Unexpected user_role, available: \"employee\" | \"moderator\"")
var ErrBadOrMissingToken = fiber.NewError(401, "The token is invalid or missing")

func NewJWT(secretKey string, AccessAllowed AccessSet) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SuccessHandler: func(c *fiber.Ctx) error {
			rawToken := c.Locals("user")
			if rawToken == nil {
				// TODO найти другой сопособ прокинуть данные об ошибке выше - не надо их вываливать на клиент
				return fiber.NewError(500, "authorize function expect token in fiber.Locals on key \"user\"")
			}

			token, ok := rawToken.(*jwt.Token)
			if !ok {
				// TODO найти другой сопособ прокинуть данные об ошибке выше - не надо их вываливать на клиент
				return fiber.NewError(500, "authorize function expect token type *jwt.Token")
			}

			claims, ok := token.Claims.(*accesscontrol.CustomClaims)
			if !ok || claims == nil {
				// TODO найти другой сопособ прокинуть данные об ошибке выше - не надо их вываливать на клиент
				return fiber.NewError(500, "authorize function expect token type *jwt.Token, with Claims type *accesscontrol.CustomClaims")
			}

			switch claims.UserRole {
			case "employee":
				if AccessAllowed.Has(AccessEmployee) {
					return c.Next()
				}
				return ErrAccessDenied
			case "moderator":
				if AccessAllowed.Has(AccessModerator) {
					return c.Next()
				}
				return ErrAccessDenied
			default:
				return ErrUnexpectedUserRole
			}
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ErrBadOrMissingToken
		},
		SigningKey: jwtware.SigningKey{Key: []byte(secretKey)},
		Claims:     &accesscontrol.CustomClaims{},
	})
}

package oauth2

import (
	"context"
	"errors"
	"fmt"
	"github.com/Simplyphotons/fyp.git/security"
	"github.com/gofiber/fiber/v2"
)

func unmatched(config *Config, c *fiber.Ctx, statusCode int) error {
	if config.allowUnmatched {
		return c.Next()
	}
	c.Response().SetStatusCode(statusCode)
	return nil
}

func (o *Config) Authorize(authorities []string) fiber.Handler {
	// Return new handler
	return func(c *fiber.Ctx) error {
		authorizationHeaders := c.GetReqHeaders()["Authorization"]
		if len(authorizationHeaders) == 0 {
			c.Response().SetStatusCode(401)
			return nil
		}
		// Get authorization header, which is passed as 'Bearer <access_token>'
		authorizationHeader := c.GetReqHeaders()["Authorization"][0]
		if authorizationHeader == "" {
			c.Response().SetStatusCode(401)
			return nil
		}

		ctx := c.UserContext()
		tokenString, err := extractToken(ctx, authorizationHeader)

		if err != nil {
			return err
		}

		authority, err := o.parseToken(ctx, tokenString)
		if err != nil {
			return err
		}

		ctx = context.WithValue(ctx, security.AuthorityKey{}, *authority)
		c.SetUserContext(ctx)

		// Validate scopes and if scopes are matching, allow request to pass through
		// Otherwise respond either with 401 or 403 code:
		//    401 - unauthenticated, usually means that there is no authorization header, or it is incorrect, for example
		//          when token is signed by the incorrect key, or the token issuer doesn't match a configured value,
		//          or audience is incorrect. The validateScopes function checks all three of these claims in the token
		//          to match preconfigured values
		//
		//    403 - when scopes in the scope claim do not match any of the configured scopes for the combination of
		//          endpoint and method
		valid, err := o.validateScopes(context.Background(), authority, authorities)
		if err != nil {
			switch err {
			case ErrInsufficientScope:
				return unmatched(o, c, 403)

			case ErrUnauthorizedRequest:
				return unmatched(o, c, 401)
			default:
				return errors.New(fmt.Sprintf("cannot validate scopes: %v", err))
			}
		} else {
			if valid {
				return c.Next()
			}
			c.Response().SetStatusCode(403)
			return nil
		}
	}
}

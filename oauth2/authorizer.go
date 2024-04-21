package oauth2

import (
	"context"
	"errors"
	"fmt"
	"github.com/Simplyphotons/fyp.git/security"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

// New creates a new instance of the OAuth2 authorization middleware
func New(config *Config) fiber.Handler {
	// Return new handler
	return func(c *fiber.Ctx) error {
		if config == nil || c.OriginalURL() == "" {
			return c.Next()
		}

		// Try to find if the requested endpoint has the authorization scopes configured
		slog.Debug("trying to find a request path match", "original_url", c.OriginalURL(), "route", c.Route().Path)
		pathRequestMatcher := config.requestMatcher[c.Route().Path]

		if pathRequestMatcher == nil {
			slog.Debug("No authorization")
			// If not, respond, depending on the middleware configuration if it allows unmatched endpoints
			// When it is configured to allow Unmatched endpoints, if middleware does not find configuration for
			// the endpoint in question, it will allow it to get handled. I.e. it means that all non-configured endpoints
			// have public access
			return unmatched(config, c, 401)
		}

		// Get scopes for the endpoint method, if not found get default endpoint scopes
		// Default endpoint scopes are applied to all methods (e.g. GET, POST, PUT etc) unless specified explicitly for the method in question
		// If there is no default scopes for the endpoint again handle it depending on the unmatched endpoints configuration
		scopes := pathRequestMatcher[c.Method()]
		if scopes == nil {
			scopes = pathRequestMatcher[""]
		}

		if scopes == nil {
			return unmatched(config, c, 401)
		}

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

		authority, err := config.parseToken(ctx, tokenString)
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
		valid, err := config.validateScopes(context.Background(), authority, scopes)
		if err != nil {
			switch err {
			case ErrInsufficientScope:
				return unmatched(config, c, 403)

			case ErrUnauthorizedRequest:
				return unmatched(config, c, 401)
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

func unmatched(config *Config, c *fiber.Ctx, statusCode int) error {
	if config.allowUnmatched {
		return c.Next()
	}
	c.Response().SetStatusCode(statusCode)
	return nil
}

package middlewares

import (
	"errors"
	"log/slog"
	"os"
	"strings"
	"ticket/api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(c *fiber.Ctx) error {
	headers := c.GetReqHeaders();

	authHeader, ok := headers[fiber.HeaderAuthorization];
	if !ok{
		return c.Status(403).JSON(fiber.Map{
			`message`: `Auth header not sent.`,
		});
	}

	if len(authHeader) != 1{
		return c.Status(403).JSON(fiber.Map{
			`message`: `Auth fail. Need token.`,
		})
	}	

	if !strings.Contains(authHeader[0], `Bearer`){
		return c.Status(403).JSON(fiber.Map{
			`message`: `Auth fail. Need Bearer token. `,
		})
	}

	bearerTokenString := strings.Split(strings.TrimSpace(authHeader[0]), ` `);

	if len(bearerTokenString) != 2{
		return c.Status(403).JSON(fiber.Map{
			`message`: `Token not set.`,
		})
	} 

	tokenString := bearerTokenString[1];

	clm := &handlers.Claims{}

	key := os.Getenv(`TOKEN_PUBLIC_KEY`);

	publicKey, err := jwt.ParseECPublicKeyFromPEM([]byte(key));
	if err != nil {
		slog.Error(`Could not parse public key from env var. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{
			`message`: `Something went wrong.`,
		})
	}

	token, err := jwt.ParseWithClaims(tokenString, clm, func(t *jwt.Token) (any, error) {

		if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New(`unexpected signing method`)
		}

		return publicKey, nil;
	})
	
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired){
			return c.Status(403).JSON(fiber.Map{
				`message`: `Token expired`,
			})
		}

		slog.Warn(`Token parse error. Error: `+ err.Error());

		return c.Status(403).JSON(fiber.Map{
			`message`: `Token auth fail.`,
		})
	}

	if !token.Valid {
		slog.Warn(`Token is not valid.`);
		return c.Status(403).JSON(fiber.Map{
			`message`: `Token auth fail.`,
		})
	}

	c.Locals(`username`, clm.Username);
	return c.Next();
}

package middlewares

import (
	"slices"

	"github.com/gofiber/fiber/v2"
)

func ForceJSONContent (c *fiber.Ctx) error{
	headers := c.GetReqHeaders();

	v, ok := headers[fiber.HeaderContentType];

	if !ok{
		return c.Status(400).JSON(fiber.Map{
			`msg`: `Content type header is not set.`,
		});
	}

	if !slices.Contains(v, fiber.MIMEApplicationJSON){
		return c.Status(400).JSON(fiber.Map{`msg`: `Content type must be application/json.`});
	}

	return c.Next();
}
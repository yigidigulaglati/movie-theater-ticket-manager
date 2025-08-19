package middlewares

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)


func CheckAdminStatus(c *fiber.Ctx) error{
	username := c.Locals(`username`);
	val, ok := username.(string);
	if !ok{	
		slog.Info(`Username value in check admin status middleware is not a string.`);
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}

	if val != `admin`{
		slog.Info(`Username is not admin in check admin status middleware.`);
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}

	return c.Next();
}




package middlewares

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)


func LogHTTP (c *fiber.Ctx) error{

	start := time.Now();

	err := c.Next();

	end := time.Now();

	duration := end.Sub(start);

	route := c.Route();
	routePath := ``
	routeName := `unnamed`;
	if route != nil{
		routePath = route.Path;
		if route.Name != ``{
			routeName = route.Name;
		}
	}else{
		routePath = c.Path();
	}

	log.Info(`REQUEST`);
	fmt.Print(`IP: `, c.IP(), `  `);
	fmt.Print(`METHOD: `, c.Method(), `  `);
	fmt.Print(`PATH: `, c.Path(), `  `);
	fmt.Print(`ROUTE PATH: `, routePath, ` `);
	fmt.Print(`ROUTE NAME: `, routeName, ` `);
	fmt.Print(`STATUS: `, c.Response().StatusCode(), `  `);
	fmt.Print(`DURATION: `, duration, `  `);

	return err;
}
package main

import (
	"gorr/api/dbInstance"
	"gorr/api/handlers"
	"gorr/api/middlewares"
	"log"

	_ "github.com/glebarez/go-sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)


func main(){

	err := godotenv.Load(`./.env`);
	if err != nil {
		log.Fatal(`Could not load env vars.`);
	}

	dbInstance.CreateDB();

	app := fiber.New();

	app.Use(middlewares.LogHTTP);
	app.Use(middlewares.ForceJSONContent);
	
	app.Post( `/login`, handlers.Login);
	app.Post(`/signup`, handlers.Signup);

	

	err = app.Listen(`:8000`);
	if err != nil {
		log.Fatal(`Could not start server. Error : ` + err.Error());
	}
}




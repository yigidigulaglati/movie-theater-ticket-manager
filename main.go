package main

import (
	"log"
	"ticket/api/dbInstance"
	"ticket/api/handlers"
	"ticket/api/middlewares"
	"ticket/api/utilities"

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

	go utilities.DeleteStaleMovieSchedules();

	app := fiber.New();

	app.Use(middlewares.LogHTTP);
	app.Use(middlewares.ForceJSONContent);

	app.Use(`/core`, middlewares.VerifyToken);
	app.Use(`/core/admin`, middlewares.CheckAdminStatus);
	
	app.Post( `/login`, handlers.Login);
	app.Post(`/signup`, handlers.Signup);

	
	app.Get(`/core/user/ticket`, handlers.BuyTicket);
	app.Get(`/core/user/movie_room`, handlers.SelectMovieRoom);
	app.Get(`/core/user/movie_rooms`, handlers.SelectMovieRooms);
	app.Get(`/core/user/movie`, handlers.SelectMovies);



	app.Post(`/core/admin/delete_movie`, handlers.DeleteMovie);
	app.Post(`/core/admin/delete_room`, handlers.DeleteRoom);
	app.Post(`/core/admin/delete_room_seating`, handlers.DeleteRoomSeating);
	app.Post(`/core/admin/insert_movie`, handlers.InsertNewMovie);
	app.Post(`/core/admin/insert_room`, handlers.InsertNewRoom);
	app.Post(`/core/admin/insert_room_seating`, handlers.InsertRoomSeating);
	app.Post(`/core/admin/schedule_movie`, handlers.ScheduleNewMovie);
	app.Post(`/core/admin/update_movie`, handlers.UpdateMovie);
	app.Post(`/core/admin/update_room`, handlers.UpdateRoom);

	err = app.Listen(`:8000`);
	if err != nil {
		log.Fatal(`Could not start server. Error : ` + err.Error());
	}
}




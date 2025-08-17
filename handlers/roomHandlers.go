package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gorr/api/dbInstance"
	"gorr/api/queries"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type Room struct{
	ID *int;
	RoomName *string;
	SeatPrice *int;
}

func InsertNewRoom(c *fiber.Ctx) error{
	room := Room{};

	err := json.Unmarshal(c.Body(), &room);
	if err != nil {
		slog.Error(`Could not decode body in insert new room handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong`})
	}

	if room.RoomName == nil || room.SeatPrice == nil{
		return c.Status(400).JSON(fiber.Map{`message`:`Must set room name and seat price of the room.`});
	}

	store := dbInstance.Store;
	
	_, err = store.Queries.InsertNewRoom(c.Context(),  queries.InsertNewRoomParams{
		RoomName: *room.RoomName,
		SeatPrice: int64(*room.SeatPrice),
	});

	if err != nil {
		slog.Error(`Could not insert new room in insert new room handler. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}
	

	return c.Status(200).JSON(fiber.Map{`message`:`Created new room.`});
}

func DeleteRoom(c *fiber.Ctx) error{
	room := Room{};
	err := json.Unmarshal(c.Body(), &room)
	if err != nil {
		slog.Error(`Could not decode body in delete room handler. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}

	if room.ID == nil{
		return c.Status(500).JSON(fiber.Map{`message`:`Must specify movie.`});
	}

	store := dbInstance.Store;

	_, err = store.Queries.DeleteRoomWithID(c.Context(), int64(*room.ID));
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return c.Status(400).JSON(fiber.Map{`message`:`No room found to delete.`});
		}
		slog.Error(`Could not delete room in delete room handler. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}


	return c.Status(200).JSON(fiber.Map{`message`:`Deleted room.`});
}

func UpdateRoom(c *fiber.Ctx) error{
	room := Room{};
	err := json.Unmarshal(c.Body(), &room)
	if err != nil {
		slog.Error(`Could not decode body in update room handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}

	if room.ID == nil || room.RoomName == nil || room.SeatPrice == nil{
		return c.Status(400).JSON(fiber.Map{`message`:`Must set room name and seat price to update room information.`});
	}

	store := dbInstance.Store;

	_, err = store.Queries.UpdateRoomWithID(c.Context(),  queries.UpdateRoomWithIDParams{
		ID: int64(*room.ID),
		SeatPrice: int64(*room.SeatPrice),
		RoomName: *room.RoomName,
	});
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return c.Status(400).JSON(fiber.Map{`message`:`No room found.`});
		}
		slog.Error(`Could not update room in update room handler. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}

	return c.Status(200).JSON(fiber.Map{`message`:`Updated room information.`});
}




























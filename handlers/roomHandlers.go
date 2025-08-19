package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"ticket/api/dbInstance"
	"ticket/api/queries"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Room struct {
	ID        *int
	RoomName  *string
	SeatPrice *int
}

type SelectMovieRoomBody struct {
	RoomMovieID *int `json:"roomMovieID"`
}


type SelectMovieRoomResponseJSONPayload struct {
	ID                 int64  `json:"id"`
	RoomMovieSeatingID int64  `json:"roomMovieSeatingID"`
	RoomID             int64  `json:"roomID"`
	MovieID            int64  `json:"movieID"`
	StartDate          string `json:"startDate"`
	RoomName           string `json:"roomName"`
	SeatPrice          int64  `json:"seatPrice"`
	Seat               int64  `json:"seat"`
	RowIndex           int64  `json:"rowIndex"`
	ColIndex           int64  `json:"colIndex"`
}

type SelectMovieRoomsBody struct {
	MovieID *int    `json:"movieID"`
	Time    *string `json:"time"`
}

type SelectRoomMovieJoinRoomJoinSeatingJSONPayload struct {
	ID            int64  `json:"id"`
	RoomID        int64  `json:"roomID"`
	MovieID       int64  `json:"movieID"`
	StartDate     string `json:"startDate"`
	RoomName      string `json:"roomName"`
	SeatPrice     int64  `json:"seatPrice"`
	UnbookedSeats int64  `json:"unbookedSeats"`
	BookedSeats   int64  `json:"bookedSeats"`
}

func InsertNewRoom(c *fiber.Ctx) error {
	room := Room{}

	err := json.Unmarshal(c.Body(), &room)
	if err != nil {
		slog.Error(`Could not decode body in insert new room handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong`})
	}

	if room.RoomName == nil || room.SeatPrice == nil {
		return c.Status(400).JSON(fiber.Map{`message`: `Must set room name and seat price of the room.`})
	}

	store := dbInstance.Store

	_, err = store.Queries.InsertNewRoom(c.Context(), queries.InsertNewRoomParams{
		RoomName:  *room.RoomName,
		SeatPrice: int64(*room.SeatPrice),
	})

	if err != nil {
		slog.Error(`Could not insert new room in insert new room handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	return c.Status(200).JSON(fiber.Map{`message`: `Created new room.`})
}

func DeleteRoom(c *fiber.Ctx) error {
	room := Room{}
	err := json.Unmarshal(c.Body(), &room)
	if err != nil {
		slog.Error(`Could not decode body in delete room handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	if room.ID == nil {
		return c.Status(500).JSON(fiber.Map{`message`: `Must specify movie.`})
	}

	store := dbInstance.Store

	_, err = store.Queries.DeleteRoomWithID(c.Context(), int64(*room.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(400).JSON(fiber.Map{`message`: `No room found to delete.`})
		}
		slog.Error(`Could not delete room in delete room handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	return c.Status(200).JSON(fiber.Map{`message`: `Deleted room.`})
}

func UpdateRoom(c *fiber.Ctx) error {
	room := Room{}
	err := json.Unmarshal(c.Body(), &room)
	if err != nil {
		slog.Error(`Could not decode body in update room handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	if room.ID == nil || room.RoomName == nil || room.SeatPrice == nil {
		return c.Status(400).JSON(fiber.Map{`message`: `Must set room name and seat price to update room information.`})
	}

	store := dbInstance.Store

	_, err = store.Queries.UpdateRoomWithID(c.Context(), queries.UpdateRoomWithIDParams{
		ID:        int64(*room.ID),
		SeatPrice: int64(*room.SeatPrice),
		RoomName:  *room.RoomName,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(400).JSON(fiber.Map{`message`: `No room found.`})
		}
		slog.Error(`Could not update room in update room handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	return c.Status(200).JSON(fiber.Map{`message`: `Updated room information.`})
}


func SelectMovieRooms(c *fiber.Ctx) error {
	b := c.Body()
	body := SelectMovieRoomsBody{}
	err := json.Unmarshal(b, &body)
	if err != nil {
		slog.Error(`Could not unmarshal body in select movie rooms handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	if body.MovieID == nil || body.Time == nil {
		slog.Info(`Nil movie id or nil limit or nil offset or nil time received in select movie rooms handler.`)
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	store := dbInstance.Store

	filterTimeUnix, err := time.Parse(time.DateTime, *body.Time)
	if err != nil {
		slog.Error(`Could not parse the time filter sent in select movie rooms. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	currentTimeUnix := time.Now().Unix();

	if filterTimeUnix.Unix() < currentTimeUnix{
		slog.Info(`Time given was previous to current time in select movie rooms handler`);
		return c.Status(400).JSON(fiber.Map{`message`:`Invalid time filter data.`});
	}

	res, err := store.Queries.SelectRoomsOfMovieJoinRoomJoinSeating(c.Context(), queries.SelectRoomsOfMovieJoinRoomJoinSeatingParams{
		MovieID:       int64(*body.MovieID),
		StartDateUnix: filterTimeUnix.Unix(),
	})

	if err != nil {
		slog.Error(`Could not execute join select in select movie rooms handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	jsonPayload := []SelectRoomMovieJoinRoomJoinSeatingJSONPayload{}
	var dateTimeFormat string
	var prevTime int64 = -1
	var prevRoom int64 = -1
	var newData SelectRoomMovieJoinRoomJoinSeatingJSONPayload
	resLength := len(res)
	unbookedSeats := 0
	bookedSeats := 0
	for i := 0; i < resLength; i++ {
		prevRoom = res[i].RoomID
		prevTime = res[i].StartDateUnix

		t := time.Unix(res[i].StartDateUnix, 0)
		dateTimeFormat = t.Format(time.DateTime)

		newData = SelectRoomMovieJoinRoomJoinSeatingJSONPayload{
			ID:            res[i].ID,
			RoomID:        res[i].RoomID,
			MovieID:       res[i].MovieID,
			StartDate:     dateTimeFormat,
			RoomName:      res[i].RoomName,
			SeatPrice:     res[i].SeatPrice,
			UnbookedSeats: 0,
			BookedSeats:   0,
		}

		unbookedSeats = 0
		bookedSeats = 0
		for j := i; j < resLength; j++ {
			if res[j].RoomID != prevRoom || res[j].StartDateUnix != prevTime {
				i = j
				break
			}
			if res[j].Seat < 0 {
				bookedSeats += 1
			} else if res[j].Seat > 0 {
				unbookedSeats++
			}
		}

		newData.BookedSeats = int64(bookedSeats)
		newData.UnbookedSeats = int64(unbookedSeats)

		jsonPayload = append(jsonPayload, newData)
	}

	return c.Status(200).JSON(fiber.Map{`message`: `qwe`, `rooms`: jsonPayload})
}


func SelectMovieRoom(c *fiber.Ctx) error {
	b := c.Body()
	bodyData := SelectMovieRoomBody{}
	err := json.Unmarshal(b, &bodyData)
	if err != nil {
		slog.Error(`Could parse the json body in select movie room handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	if bodyData.RoomMovieID == nil {
		slog.Info(`Nil room moive id found in select movie room body handler.`)
		return c.Status(400).JSON(fiber.Map{`message`: `Invalid data. ID missing.`})
	}

	store := dbInstance.Store

	rows, err := store.Queries.SelectRoomOfMovieJoinRoomJoinSeating(c.Context(), int64(*bodyData.RoomMovieID))
	if err != nil {
		slog.Error(`Could not execute join select room of movie in select movie room handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	jsonPayload := []SelectMovieRoomResponseJSONPayload{}
	var tempDate time.Time;
	for _, v := range rows {
		tempDate = time.Unix(v.StartDateUnix, 0)
		jsonPayload = append(jsonPayload, SelectMovieRoomResponseJSONPayload{
			ID:                 v.ID,
			RoomMovieSeatingID: v.RoomMovieSeatingID,
			RoomID:             v.RoomID,
			MovieID:            v.MovieID,
			StartDate:          tempDate.Format(time.DateTime),
			RoomName:           v.RoomName,
			SeatPrice:          v.SeatPrice,
			Seat:               v.Seat,
			RowIndex:           v.RowIndex,
			ColIndex:           v.ColIndex,
		});
	}

	return c.Status(200).JSON(fiber.Map{`message`:`Found room and seating information.`, `movieRoomInfo`: jsonPayload});
}

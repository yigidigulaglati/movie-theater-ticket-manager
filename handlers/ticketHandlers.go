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

type BuyTicketRequestBody struct {
	RoomMovieID        *int64 `json:"roomMovieID"`
	RoomMovieSeatingID *int64 `json:"roomMovieSeatingID"`
	RowIndex           *int64 `json:"rowIndex"`
	ColIndex           *int64 `json:"colIndex"`
}

func BuyTicket(c *fiber.Ctx) error {
	b := c.Body()
	bodyData := BuyTicketRequestBody{}
	err := json.Unmarshal(b, &bodyData)
	if err != nil {
		slog.Error(`Could not parse body in buy ticket handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	if bodyData.RoomMovieID == nil || bodyData.RoomMovieSeatingID == nil || bodyData.RowIndex == nil || bodyData.ColIndex == nil {
		slog.Info(`Got nil values in buy ticket handler.`)
		return c.Status(400).JSON(fiber.Map{`message`: `Invalid data.`})
	}

	store := dbInstance.Store

	tx, err := store.DB.BeginTx(c.Context(), nil)
	defer tx.Rollback()
	if err != nil {
		slog.Error(`Could not start transaction in buy ticket handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	qtx := store.Queries.WithTx(tx)

	seatInfo, err := qtx.SelectSeatInRoomMovieSeating(c.Context(), queries.SelectSeatInRoomMovieSeatingParams{
		RoomMovieID: *bodyData.RoomMovieID,
		RowIndex: *bodyData.RowIndex,
		ColIndex: *bodyData.ColIndex,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(400).JSON(fiber.Map{`message`: `Invalid data. No such seat.`})
		}
		slog.Error(`Could not execute select seat in buy ticket handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	timeNowUnix := time.Now().Unix();
	if seatInfo.StartDateUnix < timeNowUnix{
		slog.Info(`Seat selected for a schedule in the past in the buy ticket handler.`);
		return c.Status(400).JSON(fiber.Map{`message`:`Movie screening time has passed.`});
	}

	if seatInfo.Seat < 0 {
		return c.Status(400).JSON(fiber.Map{`message`: `Seat is already booked.`})
	}

	setSeatTo := -1
	switch seatInfo.Seat {
	case 1:
		setSeatTo = -1
	case 2:
		setSeatTo = -2
	}

	_, err = qtx.UpdateSeatInRoomMovieSeating(c.Context(), queries.UpdateSeatInRoomMovieSeatingParams{
		Seat:     int64(setSeatTo),
		RoomMovieID: seatInfo.RoomMovieID,
		RowIndex: seatInfo.RowIndex,
		ColIndex: seatInfo.ColIndex,
	})
	if err != nil {
		slog.Error(`Could not execute update seat in buy ticket handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	err = tx.Commit()
	if err != nil {
		slog.Error(`Could not commit transaction in buy ticket. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	return c.Status(200).JSON(fiber.Map{`message`: `Booked seat.`})
}







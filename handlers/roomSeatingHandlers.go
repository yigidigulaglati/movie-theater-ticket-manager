package handlers

import (
	"encoding/json"
	"gorr/api/dbInstance"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type RoomSeating struct{
	ID *int;
    RoomID *int;
    RowIndex *int;
    ColIndex *int;
    Seat *int;
}

func CreateBulkInsertQueryRoomSeating(arr []RoomSeating) (string, *[]any, error){
	var sb strings.Builder;
	_, err := sb.WriteString(`insert into seating(room_id, row_index, col_index, seat) values `);
	if err != nil {
		return ``,nil, err;
	}
	args := []any{};
	for i, s := range arr{
		_, err = sb.WriteString(`(?,?,?,?)`);
		if err != nil {
			return ``, nil, err;
		}
		if i < len(arr)-1{
			_, err = sb.WriteString(`,`);
			if err != nil {
				return ``,nil,  err;
			}
		}

		args = append(args, *(s.RoomID), *(s.RowIndex), *(s.ColIndex), *(s.Seat));
	}
	sb.WriteString(`;`);
	q := sb.String();

	return q, &args, nil;
}

func InsertRoomSeating(c *fiber.Ctx) error{
	roomSeating := []RoomSeating{};
	err := json.Unmarshal(c.Body(), &roomSeating);

	if err != nil {
		slog.Error(`Could not decode body in insert room seating handler. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}

	firstID := roomSeating[0].RoomID;
	for _, s := range roomSeating{
		if s.ColIndex == nil || s.RoomID == nil || s.RowIndex == nil || s.Seat == nil{
			slog.Info(`A nil prop found in insert new room seating.`);
			return c.Status(400).JSON(fiber.Map{`message`:`Missing data.`});
		}
		if *(s.RoomID) != *(firstID){
			slog.Info(`Different ids for rooms found in insert room seating handler.`);
			return c.Status(400).JSON(fiber.Map{`message`:`Invalid data.`});;
		}
	}

	store := dbInstance.Store;

	q, args, err := CreateBulkInsertQueryRoomSeating(roomSeating);
	if err != nil {
		slog.Error(`Could not generate bulk insert room seating in insert room seating handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}

	_, err = store.DB.ExecContext(c.Context(), q, (*args)...);
	if err != nil {
		slog.Error(`Execution of bulk insert failed in insert room seating handler. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong`});
	}

	return c.Status(200).JSON(fiber.Map{`message`:`Inserted room seating information.`});
}

func CreateBulkDeleteSeatingQuery(arr []RoomSeating) (string, *[]any, error){
	var sb strings.Builder;
	sb.WriteString(`delete from seating where `);
	n := len(arr);
	args := []any{};
	var err error;
	for i, s := range arr{
		if i == n - 1{
			_, err = sb.WriteString(`id = ?;`);
			if err != nil {
				return ``, nil, err;
			}
		}else{
			_, err = sb.WriteString(`id = ? or `);
			if err != nil {
				return ``, nil, err;
			}
		}

		args = append(args, *(s.ID));
	}

	q := sb.String();
	return q, &args, nil;
}

func DeleteRoomSeating(c *fiber.Ctx) error{
	b := c.Body();
	seating := []RoomSeating{};
	err := json.Unmarshal(b, &seating);
	if err != nil {
		slog.Error(`Could not unmarshal request body in deleteRoomSeating handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}

	prevRoomID := *(seating[0].RoomID);
	for _, s := range seating{
		if s.ID == nil{
			slog.Info(`Found nil iid in delete room seating handler.`);
			return c.Status(400).JSON(fiber.Map{`message`:`Invalid data.`});
		}
		
		if *(s.RoomID) != prevRoomID{
			slog.Info(`Found non matching room ids in delete room seating handler.`);
			return c.Status(400).JSON(fiber.Map{`message`:`Invalid data.`});
		}
	}

	q, args, err := CreateBulkDeleteSeatingQuery(seating);
	if err != nil {
		slog.Error(`Could not create bulk delete query in delete room seating handler. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong`});
	}

	store := dbInstance.Store;
	_, err = store.DB.ExecContext(c.Context(), q, *args...);
	if err != nil {
		slog.Error(`Could not execute delete query in delete room seating handler. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}

	return c.Status(200).JSON(fiber.Map{`message`:`Deleted room seatings.`});
}
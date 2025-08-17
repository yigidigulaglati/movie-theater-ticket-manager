package handlers

import (
	"encoding/json"
	"gorr/api/dbInstance"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RoomSelectResult struct {
	selectRoomID   int
	selectRowIndex int
	selectColIndex int
	selectSeat     int
}

type Schedule struct {
	RoomID    *int
	MovieID   *int
	StartDate *string
}

type ScheduleUnix struct {
	RoomID    int
	MovieID   int
	StartDate int
}

type RoomMovieIDs struct {
	RoomMovieID     int
	RoomMovieRoomID int
}

type RoomMovieSeating struct {
	RoomMovieID int
	RowIndex    int
	ColIndex    int
	Seat        int
}

// room_movie_id integer not null,

//     row_index integer not null,
//     col_index integer not null,
//     seat integer not null,

func CreateBulkInsertQueryRoomMovieSeating(data *[]RoomMovieSeating) (string, *[]any, error) {
	var sb strings.Builder
	_, err := sb.WriteString(`insert into room_movie_seating (room_movie_id, row_index, col_index, seat) values `)
	if err != nil {
		return ``, nil, err
	}

	args := []any{}
	n := len(*data)
	for i, d := range *data {
		_, err = sb.WriteString(`(?,?,?,?)`)
		if err != nil {
			return ``, nil, err
		}

		if i < n-1 {
			_, err = sb.WriteString(`,`)
			if err != nil {
				return ``, nil, err
			}
		}

		args = append(args, d.RoomMovieID, d.RowIndex, d.ColIndex, d.Seat)
	}

	_, err = sb.WriteString(`;`)
	if err != nil {
		return ``, nil, err
	}

	q := sb.String()

	return q, &args, nil
}

func CreateBulkSelectJoinRoomMovieWithSeating(data *[]RoomMovieIDs) (string, *[]any, error) {

	var sb strings.Builder
	_, err := sb.WriteString(`select room_movie.id, seating.row_index, seating.col_index, seating.seat from room_movie inner join seating on room_movie.room_id = seating.room_id where `)
	if err != nil {
		return ``, nil, err
	}

	var n = len(*data)
	var args = []any{}
	for i, d := range *data {
		_, err = sb.WriteString(`room_movie_id = ?`)
		if err != nil {
			return ``, nil, err
		}
		if i < n-1 {
			_, err = sb.WriteString(` or `)
			if err != nil {
				return ``, nil, err
			}
		}

		args = append(args, d.RoomMovieRoomID)
	}

	_, err = sb.WriteString(`;`)
	if err != nil {
		return ``, nil, err
	}
	q := sb.String()

	return q, &args, nil
}

func CreateBulkInsertQueryRoomMovie(data *[]ScheduleUnix) (string, []any, error) {

	var sb strings.Builder
	_, err := sb.WriteString(`insert into room_movie(room_id, movie_id, start_date_unix) values`)
	if err != nil {
		return ``, nil, err
	}

	args := []any{}
	n := len(*data)
	for i, s := range *data {
		_, err = sb.WriteString(`(?,?,?)`)
		if err != nil {
			return ``, nil, err
		}

		if i < n-1 {
			_, err = sb.WriteString(`,`)
			if err != nil {
				return ``, nil, err
			}
		}
		args = append(args, s.RoomID, s.MovieID, s.StartDate)

	}

	_, err = sb.WriteString(` returning id, room_id;`)
	if err != nil {
		return ``, nil, err
	}

	query := sb.String()

	return query, args, nil
}

func ScheduleNewMovie(c *fiber.Ctx) error {
	b := c.Body()
	schedules := []Schedule{}
	err := json.Unmarshal(b, &schedules)
	if err != nil {
		slog.Error(`Could not unmarshal body in schedule new movie handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	schedulesUnix := []ScheduleUnix{}
	for _, s := range schedules {
		if s.RoomID == nil || s.MovieID == nil || s.StartDate == nil {

			slog.Info(`Found nil props in schedule new movie handler`)
			return c.Status(400).JSON(fiber.Map{`message`: `Must set room, movie and start date.`})
		}

		startTime, err := time.Parse(time.DateTime, *s.StartDate)
		if err != nil {
			slog.Error(`Could not parse given time in schedule new movie. Error: ` + err.Error())
			return c.Status(400).JSON(fiber.Map{`message`: `Invalid date format.`})
		}

		startTimeUnix := startTime.Unix()
		schedulesUnix = append(schedulesUnix, ScheduleUnix{
			RoomID:    *s.RoomID,
			MovieID:   *s.MovieID,
			StartDate: int(startTimeUnix),
		})
	}

	q, args, err := CreateBulkInsertQueryRoomMovie(&schedulesUnix)
	if err != nil {
		slog.Error(`Could not create bulk insert for room movie in schedule new movie handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong`})
	}

	store := dbInstance.Store

	tx, err := store.DB.BeginTx(c.Context(), nil)
	if err != nil {
		slog.Error(`Could not start transaction in schedule new movie. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(c.Context(), q, args...)
	if err != nil {
		slog.Error(`Could not execute bulk insert in schedule new movie.. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	roomMovieIDsSlice := []RoomMovieIDs{}
	var roomMovieID int
	var roomMovieRoomID int
	for rows.Next() {
		err = rows.Scan(&roomMovieID, &roomMovieRoomID)
		if err != nil {
			slog.Error(`Could not scan the returned id in schedule new movie handler. Error: ` + err.Error())
			return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
		}

		roomMovieIDsSlice = append(roomMovieIDsSlice, RoomMovieIDs{RoomMovieID: roomMovieID, RoomMovieRoomID: roomMovieRoomID})
	}

	joinQuery, joinArgs, err := CreateBulkSelectJoinRoomMovieWithSeating(&roomMovieIDsSlice)
	if err != nil {
		slog.Error(`Could not create a bulk join in schedule new movie handler. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}

	joinRows, err := tx.QueryContext(c.Context(), joinQuery, *joinArgs...)
	if err != nil {
		slog.Error(`Could not execute join query in schedule new movie. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	var roomMovieSeatingRoomMovieID int
	var roomMovieSeatingRowIndex int
	var roomMovieSeatingColIndex int
	var roomMovieSeatingSeat int

	var roomMovieSeatingSlice = []RoomMovieSeating{}
	for joinRows.Next() {
		err = joinRows.Scan(&roomMovieSeatingRoomMovieID, &roomMovieSeatingRowIndex, &roomMovieSeatingColIndex, &roomMovieSeatingSeat)
		if err != nil {
			slog.Error(`Could not scan join result in schedule new movie handler. Error: ` + err.Error())
			return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
		}

		roomMovieSeatingSlice = append(roomMovieSeatingSlice, RoomMovieSeating{
			RoomMovieID: roomMovieSeatingRoomMovieID,
			RowIndex:    roomMovieSeatingRowIndex,
			ColIndex:    roomMovieSeatingColIndex,
			Seat:        roomMovieSeatingSeat,
		})
	}

	roomMovieSeatingQuery, roomMovieSeatingArgs, err  :=CreateBulkInsertQueryRoomMovieSeating(&roomMovieSeatingSlice);

	if err != nil {
		slog.Error(`Could not create bulk insert query room movie seating in schedule new movie handler. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}

	_, err = tx.ExecContext(c.Context(), roomMovieSeatingQuery, *roomMovieSeatingArgs...);
	if err != nil {
		slog.Error(`Could not execute bulk insert into room movie seating. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong`});
	}


	return c.Status(200).JSON(fiber.Map{`message`:`Created new scheduled movie screenings.`});
}

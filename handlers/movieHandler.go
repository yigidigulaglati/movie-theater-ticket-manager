package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"ticket/api/dbInstance"
	"ticket/api/queries"

	"github.com/gofiber/fiber/v2"
)

type SelectMovieHandlerBody struct {
	Offset      *int     `json:"offset"`
	Limit       *int     `json:"limit"`
	Title       *string  `json:"title"`
	Genre       *string  `json:"genre"`
	Rating      *float64 `json:"rating"`
}

type Movie struct {
	ID          *int     `json:"id"`
	Title       *string  `json:"title"`
	Director    *string  `json:"director"`
	ReleaseYear *int     `json:"releaseyear"`
	Genre       *string  `json:"genre"`
	Rating      *float64 `json:"rating"`
	Description *string  `json:"description"`
	Duration    *string  `json:"duration"`
}

func CreateMovieData(movieInfo *Movie) (string, sql.NullString, sql.NullInt64, sql.NullString, sql.NullFloat64, sql.NullString, sql.NullString) {

	var Title = *(movieInfo.Title)
	var Director sql.NullString
	var ReleaseYear sql.NullInt64
	var Genre sql.NullString
	var Rating sql.NullFloat64
	var Description sql.NullString
	var Duration sql.NullString

	if movieInfo.Director == nil {
		Director = sql.NullString{
			Valid:  false,
			String: *movieInfo.Director,
		}
	} else {
		Director = sql.NullString{
			Valid:  true,
			String: *movieInfo.Director,
		}
	}

	if movieInfo.ReleaseYear == nil {
		ReleaseYear = sql.NullInt64{
			Valid: false,
			Int64: int64(*movieInfo.ReleaseYear),
		}
	} else {
		ReleaseYear = sql.NullInt64{
			Valid: true,
			Int64: int64(*movieInfo.ReleaseYear),
		}
	}

	if movieInfo.Genre == nil {
		Genre = sql.NullString{
			Valid:  false,
			String: *movieInfo.Genre,
		}
	} else {
		Genre = sql.NullString{
			Valid:  true,
			String: *movieInfo.Genre,
		}
	}

	if movieInfo.Rating == nil {
		Rating = sql.NullFloat64{
			Valid: false,

			Float64: (*movieInfo.Rating),
		}
	} else {
		Rating = sql.NullFloat64{
			Valid:   true,
			Float64: (*movieInfo.Rating),
		}
	}

	if movieInfo.Description == nil {
		Description = sql.NullString{
			Valid:  false,
			String: *movieInfo.Description,
		}
	} else {
		Description = sql.NullString{
			Valid:  true,
			String: *movieInfo.Description,
		}
	}

	if movieInfo.Duration == nil {
		Duration = sql.NullString{
			Valid:  false,
			String: *movieInfo.Duration,
		}
	} else {
		Duration = sql.NullString{
			Valid:  true,
			String: *movieInfo.Duration,
		}
	}

	return Title, Director, ReleaseYear, Genre, Rating, Description, Duration
}

func InsertNewMovie(c *fiber.Ctx) error {
	movieInfo := Movie{}
	err := json.Unmarshal(c.Body(), &movieInfo)
	if err != nil {
		slog.Error(`Could not decode body in insert new movie handler. Error: ` + err.Error())
		return c.Status(400).JSON(fiber.Map{`message`: `Invalid data entered.`})
	}

	if movieInfo.Title == nil {
		return c.Status(400).JSON(fiber.Map{
			`message`: `Movie must have title.`,
		})
	}

	var Title, Director, ReleaseYear, Genre, Rating, Description, Duration = CreateMovieData(&movieInfo)

	store := dbInstance.Store

	_, err = store.Queries.InsertNewMovie(c.Context(), queries.InsertNewMovieParams{
		Title:       Title,
		Director:    Director,
		ReleaseYear: ReleaseYear,
		Genre:       Genre,
		Rating:      Rating,
		Description: Description,
		Duration:    Duration,
	})

	if err != nil {
		slog.Error(`Could not insert new movie. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{
			`message`: `Could not create a new movie.`,
		})
	}

	return c.Status(200).JSON(fiber.Map{`message`: `Created a new movie.`})
}

func DeleteMovie(c *fiber.Ctx) error {

	movie := Movie{}
	err := json.Unmarshal(c.Body(), &movie)
	if err != nil {
		slog.Error(`Could not decode body in delete movie handler. Error: ` + err.Error())
		return c.Status(400).JSON(fiber.Map{`message`: `Invalid data entered.`})
	}

	if movie.ID == nil {
		return c.Status(400).JSON(fiber.Map{`message`: `Must select a movie.`})
	}

	store := dbInstance.Store

	tx, err := store.DB.BeginTx(c.Context(), nil)
	if err != nil {
		slog.Error(`Could not begin transaction in delete movie handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}
	defer tx.Rollback()

	qtx := store.Queries.WithTx(tx)
	_, err = qtx.SelectMovieWithID(c.Context(), int64(*movie.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(400).JSON(fiber.Map{`message`: `No movie found.`})
		}

		slog.Error(`Could not select movie with id in delete movie handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong`})
	}

	_, err = qtx.DeleteMovieWithID(c.Context(), int64(*movie.ID))
	if err != nil {
		slog.Error(`Could not delete movie in delete movie handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	err = tx.Commit()
	if err != nil {
		slog.Error(`Commit failed in delete movie handler transaction. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}
	return c.Status(200).JSON(fiber.Map{`message`: `Deleted movie.`})
}

func UpdateMovie(c *fiber.Ctx) error {
	movie := Movie{}
	err := json.Unmarshal(c.Body(), &movie)
	if err != nil {
		slog.Error(`Could not decode body in update movie handler. Error: ` + err.Error())
		return c.Status(400).JSON(fiber.Map{`message`: `Invalid data.`})
	}

	if movie.ID == nil {
		slog.Info(`Movie id was not given in update movie handler.`)
		return c.Status(400).JSON(fiber.Map{`message`: `No movie id is missing.`})
	}

	store := dbInstance.Store

	var Title, Director, ReleaseYear, Genre, Rating, Description, Duration = CreateMovieData(&movie)

	_, err = store.Queries.UpdateMovieWithID(c.Context(), queries.UpdateMovieWithIDParams{
		Title:       Title,
		Director:    Director,
		ReleaseYear: ReleaseYear,
		Genre:       Genre,
		Rating:      Rating,
		Description: Description,
		Duration:    Duration,
		ID:          int64(*movie.ID),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(400).JSON(fiber.Map{`message`: `No movie found.`})
		}
		slog.Error(`Could not update movie in movie update handler. Error: ` + err.Error())
		return c.Status(500).JSON(fiber.Map{`message`: `Something went wrong.`})
	}

	return c.Status(200).JSON(fiber.Map{`message`: `Updated movie.`})
}


func CreateFilterMovieQuery(data *SelectMovieHandlerBody) (string, *[]any, error) {
	var sb strings.Builder
	_, err := sb.WriteString(`select id, title, director, release_year, genre, rating, description, duration from movies order by title `)
	if err != nil {
		return ``, nil, err
	}

	args := []any{};

	title := false
	genre := false
	rating := false
	if data.Title != nil {
		title = true
	}
	if data.Genre != nil {
		genre = true
	}
	if data.Rating != nil {
		rating = true
	}

	if title {
		_, err = sb.WriteString(`where title = ? `)
		if err != nil {
			return ``, nil, err
		}
		args = append(args, *data.Title);

		if genre {
			_, err = sb.WriteString(`and genre = ? `)
			args = append(args, *data.Genre);
			if err != nil {
				return ``, nil, err
			}
		}
		if rating {
			_, err = sb.WriteString(`and rating > ? `)
			if err != nil {
				return ``, nil, err
			}
			args = append(args, *data.Rating);
		}
	} else if genre {
		_, err = sb.WriteString(`where genre = ? `)
		if err != nil {
			return ``, nil, err
		}
		args = append(args, *data.Genre);

		if rating {
			_, err = sb.WriteString(`and rating > ? `)
			if err != nil {
				return ``, nil, err
			}
			args = append(args, *data.Rating);

		}
	} else if rating {
		_, err = sb.WriteString(`where rating > ? `)
		if err != nil {
			return ``, nil, err
		}
		args = append(args, *data.Rating);
	}

	_, err = sb.WriteString(`limit ? offset ?;`);
	if err != nil {
		return ``, nil, err
	}
	args = append(args, *data.Limit, *data.Offset);

	q := sb.String()

	return q, &args, nil;
}

func SelectMovies(c *fiber.Ctx) error {
	b := c.Body()

	bodyData := SelectMovieHandlerBody{}
	err := json.Unmarshal(b, &bodyData);
	if err != nil {
		slog.Error(`Could not parse body in select movie handler. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`})
	}

	if bodyData.Offset == nil || bodyData.Limit == nil{
		slog.Info(`offset or limit was sent as nil inside select movies`);
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}

	if *bodyData.Limit > 20{
		slog.Info(`Limit was bigger than 20.`);
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`})
	}

	q, args, err := CreateFilterMovieQuery(&bodyData);
	if err != nil {
		slog.Error(`Could not create the filter query in select movies handler. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`});
	}

	store := dbInstance.Store;

	rows, err := store.DB.QueryContext(c.Context(), q, args);
	if err != nil {
		slog.Error(`Could not execute select query in select movies. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`})
	}

	selectResult := []Movie{};
	for rows.Next(){
		var movie Movie;
		err = rows.Scan(movie.ID, movie.Title, movie.Director, movie.ReleaseYear, movie.Genre, movie.Rating, movie.Description, movie.Duration);
		if err != nil {
			slog.Error(`Could not scan select result in select movies handler. Error: ` + err.Error());
			return c.Status(500).JSON(fiber.Map{`message`:`Something went wrong.`})
		}

		selectResult = append(selectResult, movie);
	}

	return c.Status(200).JSON(fiber.Map{`message`:`Found movies.`, `movies`: selectResult});
}

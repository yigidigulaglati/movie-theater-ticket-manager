package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"ticket/api/dbInstance"
	"ticket/api/queries"
	"time"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func ValidateUsername(uname string) bool {
	
	length := utf8.RuneCountInString(uname);
	if length  < 5 || length > 20{
		return false;
	}
	return true;
}

type User struct{
	Username *string;
	Password *string;
}

type Claims struct{
	Username string;
	jwt.RegisteredClaims;
}
	
func Signup(c *fiber.Ctx) error {
	body := c.Body()
	currUser := User{}
	err := json.Unmarshal(body, &currUser)
	
	if err != nil {
		return fiber.NewError(400, `Invalid json.`)
	}	
		
	if currUser.Password == nil || currUser.Username == nil {
		return fiber.NewError(400, `Must set password and username.`)
	}

	if !ValidateUsername(*(currUser.Username)) {
		return fiber.NewError(400, `Invalid username. Must be more than 5 chars and less than 20 chars.`)
	}
	username := *(currUser.Username)
	password := *(currUser.Password)
	
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		slog.Error(`Could not hash password. Error: ` + err.Error())
		return fiber.NewError(500, `Something went wrong.`)
	}
	
	n, err := dbInstance.Store.Queries.InsertUser(context.Background(), queries.InsertUserParams{Username: username, Password: string(hash)})

	if err != nil {
		slog.Error(`DB error: ` + err.Error())
		return fiber.NewError(500, `Something went wrong. Try again later.`);
	}
	if n != 1 {
		slog.Error(`DB error: Affected number of rows is not 1.`)
		return c.Status(500).JSON(fiber.Map{
			`message`: `Something went wrong. Try again later.`,
		})
	}

	return c.Status(200).JSON(fiber.Map{`message`: `Created a Character`})
}

func Login(c *fiber.Ctx) error{
	body := c.Body();
	u := User{};
	err := json.Unmarshal(body, &u);
	if err != nil {
		slog.Error(`Could not unmarshal body. Error: ` + err.Error());
		return c.JSON(fiber.Map{
			`message`: `Invalid body.`,
		});
	}

	if u.Username == nil || u.Password == nil{
		return c.JSON(fiber.Map{
			`message`: `Must give username and password`,
		});
	}

	dbUser, err := dbInstance.Store.Queries.SelectUserWithUsername(c.Context(), *u.Username);

	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return c.Status(400).JSON(fiber.Map{
				`message`: `No User with that username exists.`,
			});
		}else{
			slog.Error(`DB err. Could not get User with username. Error: ` + err.Error());
			return c.Status(500).JSON(fiber.Map{
				`message`: `Something went wrong.`,
			});
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(*u.Password));
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword){
			return c.Status(403).JSON(fiber.Map{
				`message`: `Wrong password`,
			});
		}else{
			return c.Status(500).JSON(fiber.Map{
				`message`: `Something went wrong.`,
			});
		}
	}
	
	cl := &Claims{
		Username: *u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			Issuer: `GORR`,
			Subject: *u.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, cl)
		
	key := os.Getenv(`TOKEN_PRIVATE_KEY`);
	privateKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(key));
	if err != nil {
		slog.Error(`Could not parse private key from env var. Error: ` + err.Error());
		return c.Status(500).JSON(fiber.Map{
			`message`: `Something went wrong.`,
		})
	}

	tokenString, err := token.SignedString(privateKey);
	if err != nil {
		slog.Error(`Could not sign token. Error: ` + err.Error());
		return fiber.NewError(500, `Something went wrong. Try again later.`); 
	}

	return c.Status(200).JSON(fiber.Map{
		`message`: `Logged in.`,
		`token`: tokenString,
	});
}
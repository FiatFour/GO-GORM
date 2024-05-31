package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

func authRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	jwtSecretKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claim := token.Claims.(jwt.MapClaims)

	fmt.Println(claim["user_id"])

	return c.Next()
}

func main() {
	// Configure your PostgreSQL database details here
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // add Logger
	})
	if err != nil {
		panic("failed to connect to database")
	}

	print(db)
	fmt.Println("Connect successful")
	// db.AutoMigrate(&Book{})
	db.Migrator().CreateTable(&User{}, &Book{})

	// Setup Fiber
	app := fiber.New()
	app.Use("/books", authRequired) // Every Book API have to pass authRequired func

	//* Book API
	app.Get("/books", func(c *fiber.Ctx) error {
		return c.JSON(getBooks(db))
	})

	app.Get("/books/:id", func(c *fiber.Ctx) error {
		bookId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(getBook(db, uint(bookId)))
	})

	app.Post("/books", func(c *fiber.Ctx) error {
		book := new(Book)

		if err := c.BodyParser(book); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err := createBook(db, book)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"message": "Create Book Successful",
		})
	})

	app.Put("/books/:id", func(c *fiber.Ctx) error {
		bookId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		book := new(Book)
		if err := c.BodyParser(book); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		book.ID = uint(bookId)

		err = updateBook(db, book)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"message": "Update Book Successful",
		})
	})

	app.Delete("/books/:id", func(c *fiber.Ctx) error {
		bookId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = deleteBook(db, uint(bookId))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"message": "Delete Book Successful",
		})
	})

	//* User API

	app.Post("/register", func(c *fiber.Ctx) error {
		user := new(User)

		if err := c.BodyParser(user); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = createUser(db, user)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"message": "Register Successful",
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		user := new(User)

		if err := c.BodyParser(user); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		token, err := loginUser(db, user)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Set cookie
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 72),
			HTTPOnly: true,
		})

		return c.JSON(fiber.Map{
			// "token": token,
			"message": "Login successful!",
		})
	})

	app.Listen("localhost:8080")

	//! Create
	// newBook := &Book{
	// 	Name:        "Fiat",
	// 	Author:      "FiatFour",
	// 	Description: "test",
	// 	Price:       140,
	// }
	// createBook(db, newBook)

	//! Get(Select) and Delete
	// currentBook := getBook(db, 1)
	// currentBook.Name = "Anfat"
	// currentBook.Price = 400
	// fmt.Println(currentBook)
	// updateBook(db, currentBook)

	//! Delete
	// deleteBook(db, 1)

	//! Search Limit 1
	// searchBook := searchBook(db, "Fiat")
	// fmt.Println(searchBook)

	//* Search more than 1
	// searchBooks := searchBooks(db, "Fiat")
	// for _, book := range searchBooks {
	// 	fmt.Println(book.ID, book.Name, book.Author, book.Price)
	// }

}

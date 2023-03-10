package main

import (
	"fiber/models"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	db_name := os.Getenv("PGDATABASE")
	db_host := os.Getenv("PGHOST")
	db_pass := os.Getenv("PGPASSWORD")
	db_port := os.Getenv("PGPORT")
	db_user := os.Getenv("PGUSER")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Los_Angeles", db_host, db_user, db_pass, db_name, db_port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if (err != nil) {
		log.Fatalln("Failed to connect to DB", err)
	} else {
		log.Println("Connected to the DB")
	}

	db.AutoMigrate(models.User{})

	DB = db
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	return port
}

type User struct {
    Name string `json:"name" xml:"name" form:"name"`
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, Railway!",
		})
	})

	app.Post("/user", func(c *fiber.Ctx) error {
		user := new(User)

		if err := c.BodyParser(user); err != nil {
            return err
        }

		newUser := User{Name: user.Name}

		result := DB.Create(&user)

		if (result.Error != nil) {
			log.Fatalln("Failed to create user")
		}

		return c.Send([]byte(newUser.Name))
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		var users []models.User

		result := DB.Find(&users)

		if (result.Error != nil) {
			log.Fatalln("Failed to query users")
		}

		return c.JSON(fiber.Map{
				"users": users,
		})
	})

	app.Listen(getPort())
}

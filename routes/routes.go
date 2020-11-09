package routes

import (
	"log"
	"user-bank-manage/controllers"

	"github.com/gofiber/fiber/v2"
)

func Headerset(c *fiber.Ctx) error {
	publickey := c.Get("Public")

	if publickey != "12345" {
		return c.Status(400).JSON("header public not working")
	}
	log.Println(publickey)
	return c.Next()
}

func Init() *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, UnderWorld !!!")
	})

	bank := app.Group("banks")

	bank.Get("/fetchAll", controllers.FetchAllBanks)
	bank.Get("/detail", controllers.DetailBanks)
	bank.Post("/insert", controllers.StoreBanks)
	bank.Put("/update", controllers.UpdateBanks)
	bank.Delete("/delete", controllers.DeleteBanks)

	bank.Post("/insertJson", controllers.StoreBanksJson)
	bank.Put("/updateJson", controllers.UpdateBanksJson)

	bank.Post("/uploadImage", controllers.UploadImage)

	bank.Get("/listObject", controllers.ListObjectStorage)
	// bank.Post("/uploadObject", controllers.UploadObjectStorage)
	bank.Post("/uploadObject", controllers.MultipartObjectStorage)

	return app
}

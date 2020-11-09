package main

import (
	"user-bank-manage/db"
	"user-bank-manage/routes"
)

func main() {
	db.Init()
	index := routes.Init()
	index.Listen(":3000")
}

// func main() {
// 	app := fiber.New()
// 	db.Init()

// 	app.Get("/", func(c *fiber.Ctx) error {
// 		public := c.Get("Public-Authorization")    // "text/plain"
// 		private := c.Get("Private-Authentication") // "text/plain"
// 		// header := c.Request().Header.Method()
// 		// fmt.Println(header)
// 		fmt.Println(public)
// 		fmt.Println(private)

// 		token := strings.Split(private, "Private ")

// 		result, err := models.getAdminToken(token[1])
// 		if err != nil {
// 			return c.Status(400).JSON(http.StatusInternalServerError)
// 		}

// 		return c.Status(200).JSON(token[1])

// 		// return c.SendString("Hello Underworld !!")
// 	})

// 	app.Get("/bankManage", controllers.FetchAllBanks)
// 	app.Get("/detailBanks", controllers.DetailBanks)
// 	app.Post("/storeBanks", controllers.StoreBanks)
// 	app.Put("/updateBanks", controllers.UpdateBanks)
// 	app.Delete("/deleteBanks", controllers.DeleteBanks)

// 	app.Listen(":3000")
// }

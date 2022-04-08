package main

import (
	"fmt"

	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
	"github.com/tonystrawberry/go-fiber-crm/database"
	"github.com/tonystrawberry/go-fiber-crm/lead"
)

func setupRoutes(app *fiber.App){
	app.Get("/api/v1/lead", lead.GetLeads)
	app.Get("/api/v1/lead/:id", lead.GetLead)
	app.Post("/api/v1/lead", lead.NewLead)
	app.Delete("/api/v1/lead/:id", lead.DeleteLead)
}

func initDatabase(){
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "leads.db")
	if err != nil {
		panic("Failed to connect to database")
	}

	fmt.Println("Connected to database")
	database.DBConn.AutoMigrate(&lead.Lead{})
	fmt.Println("Database migrated")
}

func main(){
	app := fiber.New()
	initDatabase()
	setupRoutes(app)
	app.Listen(3000)

	// https://qiita.com/Ishidall/items/8dd663de5755a15e84f2
	// A defer statement defers the execution of a function until the surrounding function returns.
	// The deferred call's arguments are evaluated immediately, but the function call is not executed until the surrounding function returns.
	defer database.DBConn.Close()
}
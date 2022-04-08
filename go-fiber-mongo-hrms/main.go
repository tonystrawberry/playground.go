package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct{
	Client *mongo.Client
	Db *mongo.Database
}

var mg MongoInstance

const dbName= "fiber-hrms"
const mongoURI = "mongodb://localhost:27017/" + dbName

type Employee struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name"`
	Salary uint `json:"salary"`
	Age uint `json:"age"`
}

func Connect() error {
	client, _ := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err := client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		log.Println("Error", err);
		return err
	}

	log.Println("Connected to MongoDB");


	mg = MongoInstance {
		Client: client,
		Db: db,
	}

	return nil
}

func main() {
	if err := Connect(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()

	app.Get("/employees", func(c *fiber.Ctx) error {
		query := bson.D{{}}
		cursor, err := mg.Db.Collection("employees").Find(c.Context(), query)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		var employees []Employee = make([]Employee, 0)

		if err := cursor.All(c.Context(), &employees); err != nil {
			return c.Status(500).SendString(err.Error()) 
		}

		return c.JSON(employees)
	})

	app.Post("/employees", func(c *fiber.Ctx) error {
		collection := mg.Db.Collection("employees")

		employee := new(Employee) 

		if err := c.BodyParser(&employee); err != nil {
			return c.Status(400).SendString(err.Error())
    }

		employee.ID = ""

		result, err := collection.InsertOne(c.Context(), employee)
		if err != nil {
			log.Fatal(err)
			return c.Status(500).SendString(err.Error())
		}
		
		createdEmployee := &Employee{}

		err = collection.FindOne(c.Context(), bson.D{{Key: "_id", Value: result.InsertedID }}).Decode(&createdEmployee)
		if err != nil {
			log.Fatal(err)
			return c.Status(500).SendString(err.Error())
		}

		return c.Status(201).JSON(createdEmployee)
	})

	app.Put("/employees/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		employeeId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Fatal(err)
			return c.Status(500).SendString(err.Error())
		}

		employee := new(Employee) 
		if err := c.BodyParser(&employee); err != nil {
			return c.Status(400).SendString(err.Error())
    }


		query := bson.D{{Key: "_id", Value: employeeId }}
		updateQuery := bson.D{
			{"$set", bson.D{{"name", employee.Name}, {"salary", employee.Salary}, {"age", employee.Age}}},
    }

		collection := mg.Db.Collection("employees")

		err = collection.FindOneAndUpdate(c.Context(), query, updateQuery).Err()

		if err != nil {
			log.Fatal(err)
			if err == mongo.ErrNoDocuments {
				return c.Status(404).SendString(err.Error())
			}
			return c.Status(500).SendString(err.Error())
		}

		employee.ID = id
		return c.Status(200).JSON(employee)

	})

	app.Delete("/employees/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		employeeId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Fatal(err)
			return c.Status(500).SendString(err.Error())
		}

		query := bson.D{{Key: "_id", Value: employeeId }}

		collection := mg.Db.Collection("employees")
		result, err := collection.DeleteOne(c.Context(), &query)

		if err != nil {
			log.Fatal(err)
			return c.Status(500).SendString(err.Error())
		}

		if result.DeletedCount < 1 {
			return c.Status(404).SendString(err.Error())
		}

		return c.Status(200).SendString("Record Deleted")
	})

	log.Fatal(app.Listen(":3000"))
}
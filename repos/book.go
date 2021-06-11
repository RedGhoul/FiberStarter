package repos

import (
	"fiberstarter/database"
	"fiberstarter/models"

	"github.com/gofiber/fiber/v2"
)

func GetBooks(context *fiber.Ctx) {
	db := database.DBConn
	var books []models.Book
	//"&" generates a pointer
	//Find fill in that book array
	db.Find(&books)
	context.JSON(books)
}

func GetBook(context *fiber.Ctx) {
	id := context.Params("id")
	db := database.DBConn
	var book models.Book
	db.Find(&book, id)
	context.JSON(book)
}

func NewBook(context *fiber.Ctx) {
	params := new(struct {
		Title  string
		Author string
		Rating int
	})

	context.BodyParser(params)

	if len(params.Title) == 0 || len(params.Author) == 0 {
		context.Status(400).JSON(fiber.Map{
			"ok":    false,
			"error": "Title or Author not specified.",
		})
		return
	}

	newBook := models.Book{
		Title:  params.Title,
		Author: params.Author,
		Rating: params.Rating,
	}

	db := database.DBConn

	db.Create(&newBook)

	context.JSON(newBook)
}

func DeleteBooks(context *fiber.Ctx) {
	id := context.Params("id")
	db := database.DBConn

	var book models.Book

	db.First(&book, id)
	if book.Title == "" {
		context.Status(500).SendString("Not Book Found with given ID")
		return
	}

	db.Delete(&book)

	context.SendString("Book Deleted")

}

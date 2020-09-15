package models

import (
	"github.com/RedGhoul/bookshelf/database"
	"github.com/gofiber/fiber"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title  string `josn:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

func GetBooks(context *fiber.Ctx) {
	db := database.DBConn
	var books []Book
	//"&" generates a pointer
	//Find fill in that book array
	db.Find(&books)
	context.JSON(books)
}

func GetBook(context *fiber.Ctx) {
	id := context.Params("id")
	db := database.DBConn
	var book Book
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

	newBook := Book{
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

	var book Book

	db.First(&book, id)
	if book.Title == "" {
		context.Status(500).Send("Not Book Found with given ID")
		return
	}

	db.Delete(&book)

	context.Send("Book Deleted")

}

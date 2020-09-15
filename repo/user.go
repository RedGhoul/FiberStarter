package repos

import (
	"fmt"

	"github.com/RedGhoul/bookshelf/database"
	"github.com/RedGhoul/bookshelf/models"
	"github.com/gofiber/fiber"
)

func GetUsers(context *fiber.Ctx) {
	db := database.DBConn
	var books []models.Book
	//"&" generates a pointer
	//Find fill in that book array
	db.Find(&books)
	context.JSON(books)
}

func GetUserById(context *fiber.Ctx) {
	id := context.Params("id")
	db := database.DBConn
	var book models.Book
	db.Find(&book, id)
	context.JSON(book)
}

func GetUserByUsername(username string) models.User {
	db := database.DBConn
	var curUser models.User
	fmt.Println(username)
	db.Where("username = ?", username).First(&curUser)
	return curUser
}

func CreateUser(newUser *models.User) {
	db := database.DBConn
	db.Create(newUser)
}

func DeleteUser(context *fiber.Ctx) {
	id := context.Params("id")
	db := database.DBConn

	var book models.Book

	db.First(&book, id)
	if book.Title == "" {
		context.Status(500).Send("Not Book Found with given ID")
		return
	}

	db.Delete(&book)

	context.Send("Book Deleted")

}

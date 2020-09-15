package models

import (
	"fmt"

	"github.com/RedGhoul/bookshelf/database"
	"github.com/gofiber/fiber"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetUsers(context *fiber.Ctx) {
	db := database.DBConn
	var books []Book
	//"&" generates a pointer
	//Find fill in that book array
	db.Find(&books)
	context.JSON(books)
}

func GetUserById(context *fiber.Ctx) {
	id := context.Params("id")
	db := database.DBConn
	var book Book
	db.Find(&book, id)
	context.JSON(book)
}

func GetUserByUsername(username string) User {
	db := database.DBConn
	var curUser User
	fmt.Println(username)
	db.Where("username = ?", username).First(&curUser)
	return curUser
}

func CreateUser(newUser *User) {
	db := database.DBConn
	db.Create(newUser)
}

func DeleteUser(context *fiber.Ctx) {
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

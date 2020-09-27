package repos

import (
	"fmt"

	"github.com/RedGhoul/fiberstarter/database"
	"github.com/RedGhoul/fiberstarter/models"
)

func GetUsers() []models.User {
	db := database.DBConn
	var users []models.User
	//"&" generates a pointer
	//Find fill in that book array
	db.Find(&users)
	return users
}

func GetUserByID(userId int) models.User {
	db := database.DBConn
	var curUser models.User
	db.Find(&curUser, userId)
	return curUser
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

func DeleteUser(userId int) bool {
	db := database.DBConn

	var curUser models.User

	db.First(&curUser, userId)
	if curUser.ID == uint(userId) {
		return false
	}

	db.Delete(&curUser)
	return true
}

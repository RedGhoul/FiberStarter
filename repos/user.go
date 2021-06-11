package repos

import (
	"fmt"
	"log"

	"fiberstarter/database"
	"fiberstarter/models"
	"fiberstarter/providers"
)

func GetAllUsers() []models.User {
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

func CheckIfUserExists(username string) bool {
	db := database.DBConn
	var curUser models.User
	db.Where("username = ?", username).First(&curUser)
	if curUser.ID == 0 {
		return false
	}
	return true
}

func CreateUser(username string, email string, password string) bool {
	db := database.DBConn
	newHash, err := providers.HashProvider().CreateHash(password)
	if err != nil {
		log.Println("failed to create user")
		return false
	}
	var newUser models.User
	newUser.Email = username
	newUser.Username = username
	newUser.Password = newHash
	db.Create(&newUser)
	return true
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

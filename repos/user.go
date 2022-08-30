package repos

import (
	"log"

	"FiberStarter/database"
	"FiberStarter/models"
	"FiberStarter/providers"
)

func GetAllUsers() []models.User {
	var users []models.User
	database.DBConn.Find(&users)
	return users
}

func GetUserByID(userId int) models.User {
	var curUser models.User
	database.DBConn.Find(&curUser, userId)
	return curUser
}

func GetUserByUsername(username string) models.User {
	var curUser models.User
	database.DBConn.Where("username = ?", username).First(&curUser)
	return curUser
}

func CheckIfUserExists(email string) bool {
	var curUser models.User
	database.DBConn.Where("email = ?", email).First(&curUser)
	return curUser.ID != 0
}

func CreateUser(username string, email string, password string) bool {
	newHash, err := providers.HashProvider().CreateHash(password)
	if err != nil {
		log.Println("failed to create user")
		return false
	}
	var newUser models.User
	newUser.Email = email
	newUser.Username = username
	newUser.Password = newHash
	database.DBConn.Create(&newUser)
	return newUser.ID > 0
}

func DeleteUser(userId int) bool {
	var curUser models.User

	database.DBConn.First(&curUser, userId)
	if curUser.ID == uint(userId) {
		return false
	}

	database.DBConn.Delete(&curUser)
	return true
}

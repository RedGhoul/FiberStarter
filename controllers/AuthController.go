package controllers

import (
	"fmt"
	"log"

	"github.com/RedGhoul/bookshelf/models"
	"github.com/RedGhoul/bookshelf/providers"
	"github.com/gofiber/fiber"
)

func ShowRegisterForm(c *fiber.Ctx) {
	if err := c.Render("Auth/register", fiber.Map{}); err != nil {
		c.Status(500).Send(err.Error())
	}
}

func PostRegisterForm(c *fiber.Ctx) {
	username := c.FormValue("username")
	password1 := c.FormValue("password")
	password2 := c.FormValue("password2")
	if password1 != password2 {
		c.Send("Your passwords didn't match")
	}
	// Find user
	user := models.GetUserByUsername(username)
	if user.ID == 0 {
		fmt.Println("This is the user submitted password", password1)
		newHash, err := providers.HashProvider().CreateHash(password1)
		if err != nil {
			c.Send("Could not hash this properly")
		}
		var newUser models.User

		newUser.Email = username
		newUser.Username = username
		newUser.Password = newHash

		models.CreateUser(&newUser)
		c.Redirect("/Login")
	}
	c.Send("Could not register")
}

func ShowLoginForm(c *fiber.Ctx) {
	if err := c.Render("Auth/login", fiber.Map{}); err != nil {
		c.Status(500).Send(err.Error())
	}
}

func PostLoginForm(c *fiber.Ctx) {
	username := c.FormValue("username")
	fmt.Println(username)
	// Find user
	user := models.GetUserByUsername(username)
	fmt.Println(user)
	fmt.Println(user.Username)
	if providers.HashProvider() != nil {
		password := c.FormValue("password")
		match, err := providers.HashProvider().MatchHash(password, user.Password)
		if err != nil {
			log.Fatalf("Error when matching hash for password: %v", err)
		}
		fmt.Println(match)
		if match {
			store := providers.SessionProvider().Get(c)
			defer store.Save()
			// Set the user ID in the session store
			store.Set("userid", user.ID)
			fmt.Printf("User set in session store with ID: %v\n", user.ID)
			c.Send("You should be logged in successfully!")
		} else {
			c.Send("The entered details do not match our records.")
		}
	} else {
		panic("Hash provider was not set")
	}
}

func PostLogoutForm(c *fiber.Ctx) {
	if providers.IsAuthenticated(c) {
		store := providers.SessionProvider().Get(c)
		fmt.Println(store.Get("userid"))
		store.Delete("userid")
		store.Save()
		c.ClearCookie()
		c.Send("You are now logged out.")
	}
	c.Redirect(string(c.Fasthttp.Request.Header.Referer()))
}

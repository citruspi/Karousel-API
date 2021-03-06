package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/karousel/karousel/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func GetUserInstance(c *gin.Context) {
	db := c.MustGet("db").(gorm.DB)
	consumer := c.MustGet("consumer").(models.User)

	id := c.Params.ByName("id")

	var user models.User

	db.First(&user, id)

	if user.Username == "" {
		response := make(map[string]string)
		response["error"] = "Resource not found."
		c.JSON(404, response)
	} else {
		if user.Id != consumer.Id {
			if user.Gravatar == "" {
				user.Gravatar = user.Email
			}
			user.Email = ""
		}
		user.Password = ""
		c.JSON(200, user)
	}
}

func DeleteUserInstance(c *gin.Context) {
	db := c.MustGet("db").(gorm.DB)
	consumer := c.MustGet("consumer").(models.User)

	id := c.Params.ByName("id")

	var user models.User

	db.First(&user, id)

	if user.Username == "" {
		response := make(map[string]string)
		response["error"] = "Resource not found."
		c.JSON(404, response)
	} else {
		if (consumer.Admin) || (user.Id == consumer.Id) {
			db.Delete(&user)
			if user.Id != consumer.Id {
				if user.Gravatar == "" {
					user.Gravatar = user.Email
				}
				user.Email = ""
			}
			user.Password = ""
			c.JSON(200, user)
		} else {
			response := make(map[string]string)
			response["error"] = "Invalid credentials."
			c.JSON(401, response)
		}
	}
}

func PostUserResource(c *gin.Context) {
	db := c.MustGet("db").(gorm.DB)

	var user models.User

	c.Bind(&user)

	if (user.Username == "") || (user.Email == "") || (user.Password == "") {
		response := make(map[string]string)
		response["error"] = "Incomplete submission."
		c.JSON(400, response)
	} else {
		var queryUser models.User

		db.Where("username = ?", user.Username).First(&queryUser)

		if queryUser.Username != "" {
			response := make(map[string]string)
			response["error"] = "Duplicate resource."
			c.JSON(409, response)
		} else {
			db.Where("email = ?", user.Email).First(&queryUser)

			if queryUser.Username != "" {
				response := make(map[string]string)
				response["error"] = "Duplicate resource."
				c.JSON(409, response)
			} else {
				user.Joined = time.Now().UTC()

				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

				if err != nil {
					log.Fatal(err)
				}

				user.Password = string(hashedPassword)

				db.Create(&user)

				user.Password = ""

				locationHeader := fmt.Sprintf("/users/%v", user.Id)

				c.Writer.Header().Set("Location", locationHeader)
			}
		}
	}
}

func GetUserResource(c *gin.Context) {
	db := c.MustGet("db").(gorm.DB)

	var users []models.User
	db.Find(&users)

	for index, _ := range users {
		users[index].Password = ""
		if users[index].Gravatar == "" {
			users[index].Gravatar = users[index].Email
		}
		users[index].Email = ""
	}

	c.JSON(200, users)
}

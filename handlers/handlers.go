package handlers

import (
	"fmt"
	"log"
	"net/http"
	"pr01/db"
	"pr01/models"
	"pr01/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var login models.Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	verify := db.DB.Where("email=?", login.Email).First(&user)
	if verify.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": verify.Error.Error(),
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := utils.Jwt(login.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func SignUp(c *gin.Context) {

	var user models.SignUp
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error1": err.Error(),
		})
		return
	}

	var count int64
	db.DB.Model(&models.User{}).Where("email=?", user.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user's email already exist",
		})
		return
	}

	hashpPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {

		log.Println("error2 : ", err.Error())
		return
	}
	var user2 models.User

	user2.Email = user.Email
	user2.Name = user.Name
	user2.Password = string(hashpPass)

	result := db.DB.Create(&user2)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error3": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user2)

}

func CreateNotes(c *gin.Context) {

	email, exist := c.Get("email")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "token is missing or invalid"})
		return
	}

	var user models.User
	verify := db.DB.Where("email=?", email).First(&user)
	if verify.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "email is invalid ",
		})
		return
	}

	var notes models.Notes
	err := c.ShouldBindJSON(&notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	notes.UserId = user.UserId

	result := db.DB.Create(&notes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "note added",
	})

}

func GetNotes(c *gin.Context) {
	email, exist := c.Get("email")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized user",
		})
		return
	}

	var user models.User
	verify := db.DB.Preload("Notes").Where("email=?", email).First(&user)
	if verify.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user isnt found or something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": user,
	})

}

func DeleteNotes(c *gin.Context) {
	email, exist := c.Get("email")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "please relogin",
		})
		return
	}

	var user models.User
	verify := db.DB.Where("email=?", email).First(&user)
	if verify.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": verify.Error.Error()})
		return
	}

	id := c.Param("id")
	del := db.DB.Where("notes_id=? AND user_id=?", id, user.UserId).Delete(&models.Notes{})
	if del.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": del.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "notes deleted",
	})

}

func UpdateNotes(c *gin.Context) {

	//create a model for update
	type Update struct {
		Title *string `json:"title"`
		Body  *string `json:"body"`
	}

	//recieving the id from params (note id)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//recieving the  body from postman
	var update Update
	err = c.ShouldBindJSON(&update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println("update : ", update)
	//recieving the user's details but first recieve the email form middleware or jwt token
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	//find user according to the email
	var user models.User
	resp := db.DB.Where("email=?", email).First(&user)

	if resp.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": resp.Error.Error(),
		})
		return
	}

	if resp.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized user",
		})
		return
	}

	//fetch the note according to the id of note
	var Note models.Notes
	verify := db.DB.Model(models.Notes{}).Where("notes_id=? AND user_id=?", id, user.UserId).First(&Note)

	if verify.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": verify.Error.Error(),
		})
		return
	}

	if verify.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	//checking user input
	fmt.Println("title : ", update.Title)
	if update.Title != nil {
		Note.Notes_Title = *update.Title
		fmt.Println("title : ", Note.Notes_Title)
	}

	fmt.Println("title : ", update.Body)
	if update.Body != nil {
		Note.Notes_Body = *update.Body
		fmt.Println("Body : ", Note.Notes_Body)

	}

	r := db.DB.Save(&Note)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant update the note",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": " update successfully",
	})

}

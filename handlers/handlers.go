package handlers

import (
	"log"
	"net/http"
	"pr01/db"
	"pr01/models"
	"pr01/utils"

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
			"error": err.Error(),
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

		log.Println("error : ", err.Error())
		return
	}

	user.Password = string(hashpPass)

	result := db.DB.Create(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)

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
	del := db.DB.Where("id=? AND user_id=?", id, user.UserId).Delete(&models.Notes{})
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

func ReadNotes(c *gin.Context) {
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

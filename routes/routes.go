package routes

import (
	"pr01/handlers"
	"pr01/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {

	r.POST("/login", handlers.Login)
	r.POST("/sign-up", handlers.SignUp)

	//private
	notes := r.Group("/notes", middleware.Middleware())
	notes.POST("/create-notes", handlers.CreateNotes)
	notes.GET("/get-notes", handlers.GetNotes)
	notes.PUT("/update/:id", handlers.UpdateNotes)
	notes.DELETE("/delete-notes/:id", handlers.DeleteNotes)

}

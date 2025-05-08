package routes

import (
	"to-do-list-app/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterTodoRoutes(r *gin.Engine) {
	todo := r.Group("/api/todos")
	{
		todo.GET("/", controllers.GetTodos)
		todo.GET("/trashed", controllers.GetTrashed)
		todo.POST("/", controllers.CreateTodo)
		todo.PUT("/:id", controllers.UpdateTodo)
		todo.PUT("/order", controllers.UpdateOrder)
		todo.DELETE("/:id", controllers.DeleteTodo)
		todo.PUT("/:id/restore", controllers.RestoreTodo)
		todo.DELETE("/:id/permanent", controllers.PermanentlyDeleteTodo)
	}
}



package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func prepareV1Routes(router *gin.Engine, db *pgxpool.Pool) {

	v1 := router.Group("/v1")

	public := v1.Group("/public")
	{
		studentController := controllers.StudentController{DB: db}
		teacherController := controllers.TeacherController{DB: db}
		public.POST("/students/login", studentController.Login)
		public.POST("/teachers/login", teacherController.Login)
	}

	students := v1.Group("/students")
	students.Use(middleware.AuthMiddleware())
	{
		//studentController := controllers.StudentController{DB: db}
		students.GET("/profile", func(ctx *gin.Context) {
			userID := ctx.GetInt("user_id")
			email := ctx.GetString("email")

			ctx.JSON(200, gin.H{
				"user_id": userID,
				"email":   email,
			})
		})
	}

}

func SetupRoutes(db *pgxpool.Pool) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())
	prepareV1Routes(router, db)
	return router
}

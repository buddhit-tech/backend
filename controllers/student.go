package controllers

import (
	"backend/config"
	"backend/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentController struct {
	DB *pgxpool.Pool
}

func (c *StudentController) Login(ctx *gin.Context) {
	var req models.StudentLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id string
	var fullName string
	var password string

	err := c.DB.QueryRow(
		context.Background(),
		"SELECT id, full_name, password FROM students WHERE email=$1 AND password=$2",
		req.Email,
		req.Password,
	).Scan(&id, &fullName, &password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// ✅ Generate JWT
	token, err := config.GenerateJWT(id, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"id": id,
		"token":   token,
		"user": gin.H{
			"id":    id,
			"name":  fullName,
			"email": req.Email,
		},
	})
}

type StudentResetPasswordController struct{
	DB *pgxpool.Pool }

	func (c *StudentResetPasswordController) ResetPassword(ctx *gin.Context){
		var req models.StudentResetPasswordRequest
		if err :=ctx.ShouldBindJSON(&req);err !=nil{
			ctx.JSON(http.StatusBadRequest,gin.H{"error":"invalid JSON"})
			return
		}

		if req.Email == "" || req.NewPassword == "" {
			ctx.JSON(http.StatusBadRequest,gin.H{"error":"email and new_password are required"})
			return
		}
		
		commandTag, err := c.DB.Exec(
			context.Background(),
			"UPDATE students SET password=$1 WHERE email=$2",
			req.NewPassword,
			req.Email,
		)


		if err!= nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"failed to update student password"})
			return
		}

		if commandTag.RowsAffected() == 0{
			ctx.JSON(http.StatusNotFound,gin.H{"error":"student email not found"})
			return
		}

		ctx.JSON(http.StatusOK,gin.H{
			"status": true,
			"message":"student password reset succesfully",
		})
	}

	//Get Student by ID

	func(c*StudentController) GetStudentByID(ctx *gin.Context){
		studentID := ctx.Param("id")

		var id, fullName, email string
		err := c.DB.QueryRow(
			ctx.Request.Context(),
			"SELECT id, full_name, email FROM students WHERE id=$1",
			studentID,
		).Scan(&id, &fullName, &email,)

		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return 
		}

		ctx.JSON(http.StatusOK, gin.H{
			"id": id,
			"name": fullName,
			"email": email,
		})

	}
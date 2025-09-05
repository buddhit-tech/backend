package controllers

import (
	"backend/config"
	"backend/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeacherController struct {
	DB *pgxpool.Pool
}

func (c *TeacherController) Login(ctx *gin.Context) {
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
		"SELECT id, full_name, password FROM teachers WHERE email=$1 AND password=$2",
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
		"token":   token,
		"user": gin.H{
			"id":    id,
			"name":  fullName,
			"email": req.Email,
		},
	})
}


type TeacherResetPasswordController struct{
	DB *pgxpool.Pool }

	func (c *TeacherResetPasswordController) ResetPassword(ctx *gin.Context){
		var req models.TeacherResetPasswordRequest
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
			"UPDATE teachers SET password=$1 WHERE email=$2",
			req.NewPassword,
			req.Email,
		)

		if err!= nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"failed to update teacher password"})
			return
		}

		if commandTag.RowsAffected() == 0{
			ctx.JSON(http.StatusNotFound,gin.H{"error":"teacher email not found"})
			return
		}

		ctx.JSON(http.StatusOK,gin.H{
			"status": true,
			"message":"teacher password reset succesfully",
		})
	}
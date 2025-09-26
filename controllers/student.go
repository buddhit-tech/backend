package controllers

import (
	"backend/config"
	"backend/handlers"
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
		"token":   token,
		"user": gin.H{
			"id":    id,
			"name":  fullName,
			"email": req.Email,
		},
	})
}

func (c *StudentController) ResetPassword(ctx *gin.Context) {
	var req models.StudentResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	userID, _ := ctx.Get("user_id")

	if req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password required"})
		return
	}

	commandTag, err := c.DB.Exec(
		context.Background(),
		"UPDATE students SET password=$1 WHERE id=$2",
		req.Password,
		userID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update student password"})
		return
	}

	if commandTag.RowsAffected() == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "student email not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "student password reset succesfully",
	})
}

func (c *StudentController) GetDetails(ctx *gin.Context) {
	// Get user_id from Gin context
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}

	studentHandler := handlers.StudentHandler{DB: c.DB}
	student, err := studentHandler.FetchStudentByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch student"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   student,
	})
	return
}

func (c *StudentController) GetChatList(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}

	studentHandler := handlers.StudentHandler{DB: c.DB}
	chatList, err := studentHandler.FetchChatList(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch student"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   chatList,
	})
	return
}

func (c *StudentController) GetChatDetailsByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}

	studentHandler := handlers.StudentHandler{DB: c.DB}
	chat, err := studentHandler.FetchChatDetailsByID(userID, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch chat details"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   chat,
	})
	return
}

func (c *StudentController) GetChatMessages(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}

	studentHandler := handlers.StudentHandler{DB: c.DB}
	chat, err := studentHandler.FetchChatDetailsByID(userID, id)
	if err != nil && chat.ID == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch chat details"})
		return
	}
	chatMessages, err := studentHandler.FetchChatMessages(chat.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch student"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   chatMessages,
	})
	return
}

func (c *StudentController) GetSCSMapping(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}

	studentHandler := handlers.StudentHandler{DB: c.DB}
	scsMapping, err := studentHandler.FetchSCSDetailsByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch chat details"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   scsMapping,
	})
	return
}

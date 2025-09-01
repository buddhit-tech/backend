package routes

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

// Teacher model
type Teacher struct {
	ID        int    `json:"id"`
	TeacherID string `json:"teacher_id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone_number"`
	School    string `json:"school"`
	DOB       string `json:"dob"`
	Image     string `json:"image"`
}

// Student model
type Student struct {
	ID        int    `json:"id"`
	StudentID string `json:"student_id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone_number"`
	School    string `json:"school"`
	TeacherID string `json:"teacher_id"`
	DOB       string `json:"dob"`
	Image     string `json:"image"`
	Class     string `json:"class"`
}

// RegisterRoutes registers all routes to the Gin engine
func RegisterRoutes(r *gin.Engine, db *sql.DB, mc *memcache.Client, otpTTL time.Duration) {

	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// GET all teachers
	r.GET("/teachers", func(c *gin.Context) {
		rows, err := db.Query(`SELECT id, teacher_id, full_name, email, phone_number, school, dob, image FROM teachers`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var teachers []Teacher
		for rows.Next() {
			var t Teacher
			if err := rows.Scan(&t.ID, &t.TeacherID, &t.FullName, &t.Email, &t.Phone, &t.School, &t.DOB, &t.Image); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			teachers = append(teachers, t)
		}

		c.JSON(http.StatusOK, teachers)
	})

	// GET all students
	r.GET("/students", func(c *gin.Context) {
		rows, err := db.Query(`SELECT id, student_id, full_name, email, phone_number, school, teacher_id, dob, image, class FROM students`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var students []Student
		for rows.Next() {
			var s Student
			if err := rows.Scan(&s.ID, &s.StudentID, &s.FullName, &s.Email, &s.Phone, &s.School, &s.TeacherID, &s.DOB, &s.Image, &s.Class); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			students = append(students, s)
		}

		c.JSON(http.StatusOK, students)
	})

	// Student login (send OTP)
	r.POST("/student/login", func(c *gin.Context) {
		var req struct {
			StudentID string `json:"student_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var s Student
		err := db.QueryRow(`SELECT id, student_id, full_name FROM students WHERE student_id=$1`, req.StudentID).
			Scan(&s.ID, &s.StudentID, &s.FullName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}

		otp := "123456"
		mc.Set(&memcache.Item{
			Key:        "otp_student_" + s.StudentID,
			Value:      []byte(otp),
			Expiration: int32(otpTTL.Seconds()),
		})

		c.JSON(http.StatusOK, gin.H{"message": "OTP sent", "student": s})
	})

	// Student OTP verify
	r.POST("/student/otp/verify", func(c *gin.Context) {
		var req struct {
			UID string `json:"uid"`
			OTP string `json:"otp"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		item, err := mc.Get("otp_student_" + req.UID)
		if err != nil || string(item.Value) != req.OTP {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid OTP"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OTP verified"})
	})

	// Teacher login (send OTP)
	r.POST("/teacher/login", func(c *gin.Context) {
		var req struct {
			TeacherID string `json:"teacher_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var t Teacher
		err := db.QueryRow(`SELECT id, teacher_id, full_name FROM teachers WHERE teacher_id=$1`, req.TeacherID).
			Scan(&t.ID, &t.TeacherID, &t.FullName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "teacher not found"})
			return
		}

		otp := "123456"
		mc.Set(&memcache.Item{
			Key:        "otp_teacher_" + t.TeacherID,
			Value:      []byte(otp),
			Expiration: int32(otpTTL.Seconds()),
		})

		c.JSON(http.StatusOK, gin.H{"message": "OTP sent", "teacher": t})
	})

	// Teacher OTP verify
	r.POST("/teacher/otp/verify", func(c *gin.Context) {
		var req struct {
			UID string `json:"uid"`
			OTP string `json:"otp"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		item, err := mc.Get("otp_teacher_" + req.UID)
		if err != nil || string(item.Value) != req.OTP {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid OTP"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OTP verified"})
	})
}

package routes

import (
	"database/sql"
	"net/http"
	"time"

	"school-auth/handlers"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB, mc *memcache.Client, otpTTL time.Duration) {
	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Students
	r.GET("/students", func(c *gin.Context) {
		rows, err := db.Query(`SELECT id, student_id, full_name, email, phone_number, school, teacher_id, dob, image, class FROM students`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var students []map[string]interface{}
		for rows.Next() {
			var s = make(map[string]interface{})
			var id string
			var studentID, fullName, email, phone, school, teacherID, dob, image, class string
			if err := rows.Scan(&id, &studentID, &fullName, &email, &phone, &school, &teacherID, &dob, &image, &class); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			s["id"] = id
			s["student_id"] = studentID
			s["full_name"] = fullName
			s["email"] = email
			s["phone_number"] = phone
			s["school"] = school
			s["teacher_id"] = teacherID
			s["dob"] = dob
			s["image"] = image
			s["class"] = class
			students = append(students, s)
		}
		c.JSON(http.StatusOK, students)
	})

	// Teachers
	r.GET("/teachers", func(c *gin.Context) {
		rows, err := db.Query(`SELECT id, teacher_id, full_name, email, phone_number, school, dob, image FROM teachers`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var teachers []map[string]interface{}
		for rows.Next() {
			var t = make(map[string]interface{})
			var id string
			var teacherID, fullName, email, phone, school, dob, image string
			if err := rows.Scan(&id, &teacherID, &fullName, &email, &phone, &school, &dob, &image); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			t["id"] = id
			t["teacher_id"] = teacherID
			t["full_name"] = fullName
			t["email"] = email
			t["phone_number"] = phone
			t["school"] = school
			t["dob"] = dob
			t["image"] = image
			teachers = append(teachers, t)
		}
		c.JSON(http.StatusOK, teachers)
	})

	// Student login OTP
	r.POST("/student/login", handlers.StudentLoginHandler(db, mc, int(otpTTL.Seconds())))
	// Student OTP verify
	r.POST("/student/otp/verify", handlers.VerifyStudentOTPHandler(mc))

	// Teacher login OTP
	r.POST("/teacher/login", handlers.TeacherLoginHandler(db, mc, int(otpTTL.Seconds())))
	// Teacher OTP verify
	r.POST("/teacher/otp/verify", handlers.VerifyTeacherOTPHandler(mc))
}

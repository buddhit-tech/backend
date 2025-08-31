package routes

import (
	"database/sql"
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB, mc *memcache.Client, cfg interface{}) {
	r.GET("/teachers", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM teachers")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var teachers []map[string]interface{}
		for rows.Next() {
			var id, teacherID, fullName, email, phone, school, dob, image string
			if err := rows.Scan(&id, &teacherID, &fullName, &email, &phone, &school, &dob, &image); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			teachers = append(teachers, map[string]interface{}{
				"id":          id,
				"teacher_id":  teacherID,
				"full_name":   fullName,
				"email":       email,
				"phone_number": phone,
				"school":      school,
				"dob":         dob,
				"image":       image,
			})
		}

		c.JSON(http.StatusOK, teachers)
	})

	r.GET("/students", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM students")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var students []map[string]interface{}
		for rows.Next() {
			var id, studentID, fullName, email, phone, school, teacherID, dob, image, class string
			if err := rows.Scan(&id, &studentID, &fullName, &email, &phone, &school, &teacherID, &dob, &image, &class); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			students = append(students, map[string]interface{}{
				"id":          id,
				"student_id":  studentID,
				"full_name":   fullName,
				"email":       email,
				"phone_number": phone,
				"school":      school,
				"teacher_id":  teacherID,
				"dob":         dob,
				"image":       image,
				"class":       class,
			})
		}

		c.JSON(http.StatusOK, students)
	})
}

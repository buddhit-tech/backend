package routes

import (
	"database/sql"
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

// Teacher model
type Teacher struct {
	ID         int    `json:"id"`
	TeacherID  string `json:"teacher_id"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone_number"`
	School     string `json:"school"`
	DOB        string `json:"dob"`
	Image      string `json:"image"`
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

func RegisterRoutes(r *gin.Engine, db *sql.DB, mc *memcache.Client, cfg interface{}) {
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
}

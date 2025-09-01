package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Student struct {
    ID          string `json:"id"`
    StudentID   string `json:"student_id"`
    FullName    string `json:"full_name"`
    Email       string `json:"email"`
    PhoneNumber string `json:"phone_number"`
    School      string `json:"school"`
    TeacherID   string `json:"teacher_id"`
    DOB         string `json:"dob"`
    Class       string `json:"class"`
    Image       string `json:"image"`
}

type Teacher struct {
    ID          string `json:"id"`
    TeacherID   string `json:"teacher_id"`
    FullName    string `json:"full_name"`
    Email       string `json:"email"`
    PhoneNumber string `json:"phone_number"`
    School      string `json:"school"`
    DOB         string `json:"dob"`
    Image       string `json:"image"`
}

func main() {
    // Connect to PostgreSQL
    db, err := sql.Open("postgres", "host=157.180.16.112 user=eduroot password=F1btaRSCY8bnHvo dbname=education sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Connect to Memcached
    mc := memcache.New("127.0.0.1:11211")

    otpTTL := 300 // 5 minutes

    r := gin.Default()

    // Enable CORS for HTML frontend
    r.Use(cors.Default())

    // Student routes
    r.POST("/student/login", StudentLoginHandler(db, mc, otpTTL))
    r.POST("/student/otp/verify", VerifyStudentOTPHandler(mc))

    // Teacher routes
    r.POST("/teacher/login", TeacherLoginHandler(db, mc, otpTTL))
    r.POST("/teacher/otp/verify", VerifyTeacherOTPHandler(mc))

    fmt.Println("Server running on http://localhost:8080")
    r.Run(":8080")
}

// -------------------- Handlers --------------------

func StudentLoginHandler(db *sql.DB, mc *memcache.Client, otpTTLSeconds int) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            StudentID string `json:"student_id"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
            return
        }

        // Fetch student from DB
        var student Student
        row := db.QueryRow(`SELECT id, student_id, full_name, email, phone_number, school, teacher_id, dob, class, image 
            FROM students WHERE student_id=$1`, req.StudentID)
        err := row.Scan(&student.ID, &student.StudentID, &student.FullName, &student.Email, &student.PhoneNumber,
            &student.School, &student.TeacherID, &student.DOB, &student.Class, &student.Image)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
            return
        }

        // Generate UID & OTP
        uid := generateUID()
        otp := generateOTP()

        // Store OTP in Memcached
        mc.Set(&memcache.Item{
            Key:        "student_otp_" + uid,
            Value:      []byte(otp),
            Expiration: int32(otpTTLSeconds),
        })

        // Print OTP in terminal
        fmt.Printf("[StudentLogin] StudentID: %s | UID: %s | OTP: %s\n", req.StudentID, uid, otp)

        // Return JSON
        c.JSON(http.StatusOK, gin.H{
            "uid":     uid,
            "otp":     otp,
            "message": "OTP sent successfully",
            "student": student,
        })
    }
}

func VerifyStudentOTPHandler(mc *memcache.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            UID string `json:"uid"`
            OTP string `json:"otp"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
            return
        }

        item, err := mc.Get("student_otp_" + req.UID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired OTP"})
            return
        }

        if string(item.Value) != req.OTP {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid OTP"})
            return
        }

        // Generate session token
        sessionToken := generateSessionToken()
        mc.Set(&memcache.Item{
            Key:        "student_session_" + req.UID,
            Value:      []byte(sessionToken),
            Expiration: int32(time.Hour.Seconds()),
        })

        c.JSON(http.StatusOK, gin.H{"session_token": sessionToken})
    }
}

func TeacherLoginHandler(db *sql.DB, mc *memcache.Client, otpTTLSeconds int) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            TeacherID string `json:"teacher_id"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
            return
        }

        var teacher Teacher
        row := db.QueryRow(`SELECT id, teacher_id, full_name, email, phone_number, school, dob, image 
            FROM teachers WHERE teacher_id=$1`, req.TeacherID)
        err := row.Scan(&teacher.ID, &teacher.TeacherID, &teacher.FullName, &teacher.Email, &teacher.PhoneNumber,
            &teacher.School, &teacher.DOB, &teacher.Image)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "teacher not found"})
            return
        }

        uid := generateUID()
        otp := generateOTP()

        mc.Set(&memcache.Item{
            Key:        "teacher_otp_" + uid,
            Value:      []byte(otp),
            Expiration: int32(otpTTLSeconds),
        })

        fmt.Printf("[TeacherLogin] TeacherID: %s | UID: %s | OTP: %s\n", req.TeacherID, uid, otp)

        c.JSON(http.StatusOK, gin.H{
            "uid":     uid,
            "otp":     otp,
            "message": "OTP sent successfully",
            "teacher": teacher,
        })
    }
}

func VerifyTeacherOTPHandler(mc *memcache.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            UID string `json:"uid"`
            OTP string `json:"otp"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
            return
        }

        item, err := mc.Get("teacher_otp_" + req.UID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired OTP"})
            return
        }

        if string(item.Value) != req.OTP {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid OTP"})
            return
        }

        sessionToken := generateSessionToken()
        mc.Set(&memcache.Item{
            Key:        "teacher_session_" + req.UID,
            Value:      []byte(sessionToken),
            Expiration: int32(time.Hour.Seconds()),
        })

        c.JSON(http.StatusOK, gin.H{"session_token": sessionToken})
    }
}

// -------------------- Helper Functions --------------------

func generateUID() string {
    b := make([]byte, 4)
    rand.Read(b)
    return hex.EncodeToString(b)
}

func generateOTP() string {
    b := make([]byte, 4)
    rand.Read(b)
    return fmt.Sprintf("%06d", int(b[0])<<16|int(b[1])<<8|int(b[2])%1000000)
}

func generateSessionToken() string {
    b := make([]byte, 16)
    rand.Read(b)
    return hex.EncodeToString(b)
}

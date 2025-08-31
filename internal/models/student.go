package models

// Student represents a student record in the system
type Student struct {
	ID        string `json:"id"`         // unique DB ID
	StudentID string `json:"student_id"` // school-specific ID
	FullName  string `json:"full_name"`  // full name of the student
	Email     string `json:"email"`      // email address
	Phone     string `json:"phone"`      // phone number
	School    string `json:"school"`     // school name
	TeacherID string `json:"teacher_id"` // associated teacher's ID
	DOB       string `json:"dob"`        // date of birth (YYYY-MM-DD)
	Image     string `json:"image"`      // URL or path to profile image
	Class     string `json:"class"`      // current class/grade
}

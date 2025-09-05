package models

type Student struct {
	ID        string `json:"id"`
	StudentID string `json:"student_id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone_number"`
	School    string `json:"school"`
	TeacherID string `json:"teacher_id"`
	DOB       string `json:"dob"`
	Image     string `json:"image"`
	Class     string `json:"class"`
	Password  string `json:"password"`
}

type StudentLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


type StudentResetPasswordRequest struct {
	Email    string `json:"email"`
	NewPassword string `json:"new_password"`
}

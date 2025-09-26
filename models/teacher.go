package models

type Teacher struct {
	ID        string `json:"id"`
	TeacherID string `json:"teacher_id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone_number"`
	School    string `json:"school"`
	DOB       string `json:"dob"`
	Image     string `json:"image"`
	Password  string `json:"password"`
}

type TeacherResetPasswordRequest struct {
	Password string `json:"password"`
}

package models

type Student struct {
	ID        string  `db:"id" json:"id"`
	StudentID string  `db:"student_id" json:"student_id"`
	FullName  string  `db:"full_name" json:"full_name"`
	Email     string  `db:"email" json:"email"`
	Phone     string  `db:"phone_number" json:"phone_number"`
	DOB       *string `db:"dob" json:"dob"`
	Image     *string `db:"image" json:"image"`
	Password  *string `db:"password" json:"password"`
}

type StudentLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type StudentResetPasswordRequest struct {
	Password string `json:"password"`
}

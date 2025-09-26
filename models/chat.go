package models

import "time"

type PublicChat struct {
	ID              string     `db:"id" json:"id"`
	StudentID       *string    `db:"student_id" json:"student_id"`
	Title           *string    `db:"title" json:"title"`
	Description     *string    `db:"description" json:"description"`
	TeacherGlobalId *string    `db:"teacher_global_id" json:"teacher_global_id"`
	TeacherId       *string    `db:"teacher_id" json:"teacher_id"`
	CreatedAt       *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updated_at"`
}

type PublicChatMessage struct {
	ID         string     `db:"id" json:"id"`
	ChatID     string     `db:"chat_id" json:"chat_id"`
	Question   string     `db:"question" json:"question"`
	Answer     *string    `db:"answer" json:"answer"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`   // when question stored
	AnsweredAt *time.Time `db:"answered_at" json:"answered_at"` // when AI responded
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}

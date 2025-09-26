package handlers

import (
	"backend/models"
	"context"
	"encoding/json"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentHandler struct {
	DB *pgxpool.Pool
}

// ------------------
// Request payloads
// ------------------
type StudentLoginRequest struct {
	StudentID string `json:"student_id"`
}

type VerifyStudentOTPRequest struct {
	UID string `json:"uid"`
	OTP string `json:"otp"`
}

// ------------------
// DB Helper
// ------------------
func (c *StudentHandler) FetchStudentByID(id string) (models.Student, error) {
	var student models.Student
	query := `SELECT id, full_name, email, phone_number, image
              FROM students WHERE id=$1`
	err := pgxscan.Get(context.Background(), c.DB, &student, query, id)
	if err != nil {
		return models.Student{}, err
	}
	return student, nil
}

func (c *StudentHandler) FetchChatList(id string) ([]models.PublicChat, error) {
	var publicChats []models.PublicChat
	query := `SELECT * FROM public_chats WHERE student_id=$1`
	err := pgxscan.Select(context.Background(), c.DB, &publicChats, query, id)
	if err != nil {
		return []models.PublicChat{}, err
	}
	return publicChats, nil
}

func (c *StudentHandler) FetchChatDetailsByID(userID string, chatId string) (models.PublicChat, error) {
	var publicChat models.PublicChat
	query := `SELECT * FROM public_chats WHERE id=$1 AND student_id=$2`
	err := pgxscan.Get(context.Background(), c.DB, &publicChat, query, chatId, userID)
	if err != nil {
		return models.PublicChat{}, err
	}
	return publicChat, nil
}

func (c *StudentHandler) FetchChatMessages(chatID string) ([]models.PublicChatMessage, error) {
	var publicMessages []models.PublicChatMessage
	query := `SELECT * FROM public_messages WHERE chat_id=$1 ORDER BY created_at DESC`
	err := pgxscan.Select(context.Background(), c.DB, &publicMessages, query, chatID)
	if err != nil {
		return []models.PublicChatMessage{}, err
	}
	return publicMessages, nil
}

func (c *StudentHandler) FetchSCSDetailsByUserID(userID string) ([]models.YearWiseDetails, error) {
	query := `
		SELECT 
			scs.year,
			s_scs.is_active,
			COALESCE(
				json_agg(
					json_build_object(
						'scs_id', s_scs.scs_id,
						'school', schools.name,
						'class', classes.name,
						'subject', subjects.name
					) ORDER BY subjects.name
				), '[]'::json
			) AS details
		FROM student_scs_mapping AS s_scs
		JOIN school_class_subject_mapping AS scs ON s_scs.scs_id = scs.id 
		JOIN schools ON schools.id = scs.school_id 
		JOIN classes ON classes.id = scs.class_id 
		JOIN subjects ON subjects.id = scs.subject_id 
		WHERE s_scs.student_id = $1
		GROUP BY scs.year, s_scs.is_active 
		ORDER BY scs.year;
	`

	rows, err := c.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.YearWiseDetails
	for rows.Next() {
		var ywd models.YearWiseDetails
		var detailsRaw []byte

		if err := rows.Scan(&ywd.Year, &ywd.IsActive, &detailsRaw); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(detailsRaw, &ywd.Details); err != nil {
			return nil, err
		}

		results = append(results, ywd)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

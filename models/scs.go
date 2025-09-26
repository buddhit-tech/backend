package models

type SCSDetail struct {
	ScsID   string `json:"scs_id"`
	School  string `json:"school"`
	Class   string `json:"class"`
	Subject string `json:"subject"`
}

type YearWiseDetails struct {
	Year     int         `json:"year"`
	IsActive bool        `json:"is_active"`
	Details  []SCSDetail `json:"details"`
}

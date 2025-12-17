package models

type ActivityLog struct {
	ActivityLogID string `json:"activity_log_id"`
	UserID        string `json:"user_id"`
	Category      string `json:"category"`
	ActivityName  string `json:"activity_name"`
	IsSuccess     bool   `json:"is_success"`
	Metadata      string `json:"metadata"`
	Description   string `json:"description"`
	Timestamp     string `json:"timestamp"`
}

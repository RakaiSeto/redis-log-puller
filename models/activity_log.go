package models

type ActivityLog struct {
	ActivityLogID string `json:"activity_log_id"`
	UserID        string `json:"user_id"`

	// "expense", "income", "auth", dll
	Category string `json:"category"`

	// "create", "update", "delete", "login"
	ActivityName string `json:"activity_name"`

	// TARGET aktivitas
	EntityType string `json:"entity_type"` // "expense"
	EntityID   string `json:"entity_id"`   // uuid expense

	IsSuccess  bool   `json:"is_success"`
	Description string `json:"description"`
	Metadata    string `json:"metadata"`

	Timestamp string `json:"timestamp"`
}

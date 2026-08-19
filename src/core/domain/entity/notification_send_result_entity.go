package entity

type NotificationSendResultEntity struct {
	ID             uint     `json:"id"`
	NotificationID uint     `json:"notification_id"`
	SuccessCount   int      `json:"success_count"`
	FailureCount   int      `json:"failure_count"`
	FailedTokens   []string `json:"failed_tokens"`
	Errors         []string `json:"errors"`
}

func NewNotificationSendResultEntity(ID uint, notificationID uint, successCount, failureCount int, failedTokens []string, errors []string) *NotificationSendResultEntity {
	return &NotificationSendResultEntity{ID: ID, NotificationID: notificationID, SuccessCount: successCount, FailureCount: failureCount, FailedTokens: failedTokens, Errors: errors}
}

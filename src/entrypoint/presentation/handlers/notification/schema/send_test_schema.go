package schema

type SendTestNotificationRequest struct {
	Topic  string   `json:"topic"`
	Tokens []string `json:"tokens"`
	Title  string   `json:"title"`
	Body   string   `json:"body"`
}

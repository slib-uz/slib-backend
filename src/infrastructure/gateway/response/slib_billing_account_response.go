package response

type SlibBillingAccountResponse struct {
	Status  int                        `json:"status"`
	Success bool                       `json:"success"`
	Data    SlibBillingAccountData     `json:"data"`
}

type SlibBillingAccountData struct {
	ID        int64  `json:"id"`
	JournalID int64  `json:"journal_id"`
	UserID    int64  `json:"user_id"`
	Type      int    `json:"type"`
	TypeName  string `json:"type_name"`
	Balance   int64  `json:"balance"`
}

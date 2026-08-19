package response

type SubmissionResubmitResponse struct {
	ID          uint   `json:"id"`
	Deadline    string `json:"deadline"`
	OldDeadline string `json:"old_deadline"`
}

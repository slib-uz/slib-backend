package enum

type ClaimStatus string

const (
	ClaimStatusPending  ClaimStatus = "pending"
	ClaimStatusApproved ClaimStatus = "approved"
	ClaimStatusRejected ClaimStatus = "rejected"
)

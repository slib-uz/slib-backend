package enum

type PeerReviewStatus int

const (
	PeerReviewStatusPending        PeerReviewStatus = 1
	PeerReviewStatusAccepted       PeerReviewStatus = 2
	PeerReviewStatusResubmitted    PeerReviewStatus = 3
	PeerReviewStatusRejected       PeerReviewStatus = -2
	PeerReviewStatusDeadlineMissed PeerReviewStatus = -1

	PeerReviewStatusMinorRevision PeerReviewStatus = 10
	PeerReviewStatusMajorRevision PeerReviewStatus = 20

	PeerReviewStatusRecommended    PeerReviewStatus = 30
	PeerReviewStatusNotRecommended PeerReviewStatus = -30
)

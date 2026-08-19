package enum

type PeerReviewMethod int

const (
	SingleBlind PeerReviewMethod = 1 // Reviewer is anonymous to author
	DoubleBlind PeerReviewMethod = 2 // Both author and reviewer are anonymous
	TripleBlind PeerReviewMethod = 3 // Author, reviewer, and editor are anonymous
	OpenReview  PeerReviewMethod = 4 // All parties are known to each other
)

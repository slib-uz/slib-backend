package enum

import "strconv"

type Stage int

const (
	TechnicalReviewStage Stage = 10
	PeerReviewStage      Stage = 20
	PaymentStage         Stage = 30
	PublishStage         Stage = 40
)

func (this Stage) ToString() string {
	return strconv.FormatInt(int64(this), 10)
}

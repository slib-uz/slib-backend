package entity

import "slib.uz/src/core/domain/entity/enum"

type CrossRefDepositResultEntity struct {
	Status       enum.DoiDepositStatus
	DOI          string
	Message      string
	SubmissionID string
	RequestBody  string
	ResponseBody string
}

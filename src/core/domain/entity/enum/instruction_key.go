package enum

type InstructionKey string

const (
	InstructionForAuthors   InstructionKey = "FOR_AUTHORS"
	InstructionForReviewers InstructionKey = "FOR_REVIEWERS"
	InstructionSubmission   InstructionKey = "SUBMISSION"
	InstructionEthics       InstructionKey = "ETHICS"
	InstructionOpenAccess   InstructionKey = "OPEN_ACCESS"
)

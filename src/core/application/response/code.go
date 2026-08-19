package response

type Code struct {
	Code int
	Name string
}

func NewCode(code int, name string) *Code {
	return &Code{Code: code, Name: name}
}

var (
	CodeJournalAlreadyExists     = NewCode(4100, "JOURNAL_ALREADY_EXISTS")
	CodeApplicationAlreadyExists = NewCode(4101, "APPLICATION_ALREADY_EXISTS")

	CodeRoleAlreadyExists = NewCode(4102, "ROLE_ALREADY_EXISTS")
	RoiNotFound           = NewCode(4103, "ROI_NOT_FOUND")
	RoiFound              = NewCode(4104, "ROI_FOUND")
	RoiUnexpectedError    = NewCode(4500, "ROI_UNEXPECTED_ERROR")

	CodeNotFound      = NewCode(4004, "NOT_FOUND")
	CodeAlreadyExists = NewCode(4005, "ALREADY_EXISTS")
	CodeNotReady      = NewCode(4006, "NOT_READY")
	CodeConflictError = NewCode(4009, "CONFLICT_ERROR")

	CodeOrcidAlreadyConnected = NewCode(4110, "ORCID_ALREADY_CONNECTED")
	CodeInvalidOrReusedCode   = NewCode(4111, "INVALID_OR_REUSED_CODE")
	CodeExpired               = NewCode(4112, "EXPIRED_CODE")

	CodeInsufficientBalance       = NewCode(4100, "INSUFFICIENT_BALANCE")
	CodeEditionHasArticles        = NewCode(4105, "EDITION_HAS_ARTICLES")
	CodePublisherTinAlreadyExists = NewCode(4106, "PUBLISHER_TIN_ALREADY_EXISTS")
	CodeOrganizationAlreadyExists = NewCode(4107, "ORGANIZATION_ALREADY_EXISTS")

	CodeCrossRefForbidden = NewCode(4200, "CROSSREF_FORBIDDEN")

	CodeIntegrationError = NewCode(4500, "INTEGRATION_ERROR")
)

package enum

type NotificationTemplate string

const (
	TemplateAcceptedAtTechnicalReview     NotificationTemplate = "APPLICATION_ACCEPTED_AT_TECHNICAL_REVIEW"
	TemplateRejectedAtTechnicalReview     NotificationTemplate = "APPLICATION_REJECTED_AT_TECHNICAL_REVIEW"
	TemplateApplicationReSent             NotificationTemplate = "APPLICATION_RESENT"
	TemplateDeadlineDateAtTechnicalReview NotificationTemplate = "DEADLINE_DATE_AT_TECHNICAL_REVIEW"
	TemplateDeadlineDateAtPeerReview      NotificationTemplate = "DEADLINE_DATE_AT_PEER_REVIEW"
	TemplateDeadlineDateAtPaymentStage    NotificationTemplate = "DEADLINE_DATE_AT_PAYMENT_STAGE"
	TemplateDeadlineDateAtPublishStage    NotificationTemplate = "DEADLINE_DATE_AT_PUBLISH_STAGE"
	TemplateNewPeerReview                 NotificationTemplate = "NEW_PEER_REVIEW"
	TemplateSpellCheckFinished            NotificationTemplate = "SPELL_CHECK_FINISHED"
	TemplateAntiPlagCheckFinished         NotificationTemplate = "ANTI_PLAGIARISM_CHECK_FINISHED"
	TemplateAiDetectCheckFinished         NotificationTemplate = "AI_DETECT_CHECK_FINISHED"
	TemplateNewApplication                NotificationTemplate = "NEW_APPLICATION"
	TemplateApplicationPublished          NotificationTemplate = "APPLICATION_PUBLISHED"
	TemplateAcceptedAtPeerReview          NotificationTemplate = "APPLICATION_ACCEPTED_AT_PEER_REVIEW"
	TemplateRejectedAtPeerReview          NotificationTemplate = "APPLICATION_REJECTED_AT_PEER_REVIEW"
	TemplatePayToPublishArticle           NotificationTemplate = "PAY_TO_PUBLISH_ARTICLE"
	TemplateApplicationPaymentSuccessful  NotificationTemplate = "APPLICATION_PAYMENT_SUCCESSFUL"
	TemplateNewSupportQuestion            NotificationTemplate = "NEW_SUPPORT_QUESTION"
	TemplateNewSupportAnswer              NotificationTemplate = "NEW_SUPPORT_ANSWER"
)
